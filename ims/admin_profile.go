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

func (i Config) validateGetAdminProfileConfig() error {
	switch i.ProfileApiVersion {
	case "v1", "v2", "v3":
	default:
		return fmt.Errorf("invalid API version parameter, latest version is v3")
	}

	switch {
	case i.ServiceToken == "":
		return fmt.Errorf("missing service token parameter")
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client ID parameter")
	case i.Guid == "":
		return fmt.Errorf("missing guid parameter")
	case i.AuthSrc == "":
		return fmt.Errorf("missing auth source parameter")

	default:
		log.Println("all needed parameters verified not empty")
	}
	return nil
}

// GetAdminProfile requests the user profile using an access token.
func (i Config) GetAdminProfile() (string, error) {

	err := i.validateGetAdminProfileConfig()
	if err != nil {
		return "", fmt.Errorf("invalid parameters for admin profile: %w", err)
	}

	httpClient, err := i.httpClient()
	if err != nil {
		return "", fmt.Errorf("error creating the HTTP Client: %w", err)
	}

	c, err := ims.NewClient(&ims.ClientConfig{
		URL:    i.URL,
		Client: httpClient,
	})
	if err != nil {
		return "", fmt.Errorf("error creating the client: %w", err)
	}

	profile, err := c.GetAdminProfile(&ims.GetAdminProfileRequest{
		ServiceToken: i.ServiceToken,
		ApiVersion:   i.ProfileApiVersion,
		ClientID:     i.ClientID,
		Guid:         i.Guid,
		AuthSrc:      i.AuthSrc,
	})
	if err != nil {
		return "", fmt.Errorf("error getting admin profile: %w", err)
	}

	return string(profile.Body), nil
}
