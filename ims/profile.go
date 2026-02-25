// Copyright 2020 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

package ims

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/adobe/ims-go/ims"
)

func (i Config) validateGetProfileConfig() error {
	switch i.ProfileApiVersion {
	case "v1", "v2", "v3":
	default:
		return fmt.Errorf("invalid API version parameter, latest version is v3")
	}

	switch {
	case i.AccessToken == "":
		return fmt.Errorf("missing access token parameter")
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	default:
		log.Println("all needed parameters verified not empty")
	}
	return nil
}

// GetProfile requests the user profile using an access token.
func (i Config) GetProfile() (string, error) {

	err := i.validateGetProfileConfig()
	if err != nil {
		return "", fmt.Errorf("invalid parameters for profile: %v", err)
	}

	httpClient, err := i.httpClient()
	if err != nil {
		return "", fmt.Errorf("error creating the HTTP Client: %v", err)
	}

	c, err := ims.NewClient(&ims.ClientConfig{
		URL:    i.URL,
		Client: httpClient,
	})
	if err != nil {
		return "", fmt.Errorf("error creating the client: %v", err)
	}

	profile, err := c.GetProfile(&ims.GetProfileRequest{
		AccessToken: i.AccessToken,
		ApiVersion:  i.ProfileApiVersion,
	})
	if err != nil {
		return "", err
	}

	if !i.DecodeFulfillableData {
		return string(profile.Body), nil
	}

	// Decode the fulfillable_data in the product context
	decodedProfile, err := decodeProfile(profile.Body)
	if err != nil {
		return "", err
	}
	return decodedProfile, nil
}

func decodeProfile(profile []byte) (string, error) {
	// Parse the profile JSON
	var p map[string]interface{}
	err := json.Unmarshal(profile, &p)
	if err != nil {
		return "", fmt.Errorf("error parsing profile JSON: %v", err)
	}
	findFulfillableData(p)

	modifiedJson, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON during profile decode: %v", err)
	}

	return string(modifiedJson), nil
}

func findFulfillableData(data interface{}) {
	switch data := data.(type) {
	case map[string]interface{}:
		for key, value := range data {
			if key == "fulfillable_data" {
				serviceCode, ok := data["serviceCode"].(string)
				if ok && (serviceCode == "dma_media_library" ||
					serviceCode == "dma_aem_cloud" ||
					serviceCode == "dma_aem_contenthub" ||
					serviceCode == "dx_genstudio") {

					decodedFulfillableData, err := modifyFulfillableData(value.(string))
					if err != nil {
						fmt.Printf("Error decoding fulfillable_data: %v", err)
						return
					}
					data["fulfillable_data"] = decodedFulfillableData
				}
			} else {
				findFulfillableData(value)
			}
		}
	case []interface{}:
		for _, item := range data {
			findFulfillableData(item)
		}
	}
}

type fulfillableData struct {
	Iid string `json:"iid"`
}

func modifyFulfillableData(data string) (string, error) {
	strippedGzippedInstanceID := strings.Replace(data, "\"", "", 2)
	gzippedInstanceIDBytes, err := base64.StdEncoding.DecodeString(strippedGzippedInstanceID)
	if err != nil {
		return "", fmt.Errorf("unable to base64 decode fulfillable_data: %v", err)
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(gzippedInstanceIDBytes))
	if err != nil {
		return "", fmt.Errorf("unable to create gzip reader: %v", err)
	}
	defer func() {
		if _, gzErr := io.Copy(io.Discard, gzipReader); gzErr != nil {
			log.Printf("error while consuming the gzip reader: %v", gzErr)
		}

		if gzErr := gzipReader.Close(); gzErr != nil {
			log.Printf("unable to close gzip reader: %v", gzErr)
		}
	}()

	iidDecoder := json.NewDecoder(gzipReader)
	instanceIdJson := fulfillableData{}
	err = iidDecoder.Decode(&instanceIdJson)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshall the fulfillable_data: %v", err)
	}
	return instanceIdJson.Iid, nil
}
