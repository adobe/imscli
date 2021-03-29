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
	"fmt"
	"log"

	"github.com/adobe/ims-go/ims"
)

func (i Config) validateGetOrganizationsConfig() error {
	switch {
	case i.AccessToken == "":
		return fmt.Errorf("missing access token parameter")
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case i.OrgsApiVersion == "":
		return fmt.Errorf("missing API version parameter")
	default:
		log.Println("all needed parameters verified not empty")
	}
	return nil
}

/*
 * GetProfile requests the user profile using an access token.
 */
func (i Config) GetOrganizations() (string, error) {

	err := i.validateGetOrganizationsConfig()
	if err != nil {
		return "", fmt.Errorf("invalid parameters for organizations: %v", err)
	}

	httpClient, err := i.httpClient()
	if err != nil {
		return "", fmt.Errorf("error creating the HTTP Client: %v", err)
	}

	c, err := ims.NewClient(&ims.ClientConfig{
		URL: i.URL,
		Client: httpClient,
	})
	if err != nil {
		return "", fmt.Errorf("error creating the client: %v", err)
	}

	organizations, err := c.GetOrganizations(&ims.GetOrganizationsRequest{
		AccessToken: i.AccessToken,
		ApiVersion: i.OrgsApiVersion,
	})
	if err != nil {
		return "", err
	}

	return string(organizations.Body), nil
}
