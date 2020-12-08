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

	"github.com/adobe/ims-go/ims"
	"github.com/telegrapher/vrb"
)

func (i Config) validateGetProfileConfig() error {
	switch {
	case i.UserToken == "":
		return fmt.Errorf("missing access token parameter")
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	default:
		vrb.Println("all needed parameters verified not empty")
	}
	return nil
}

/*
 * GetProfile requests the user profile using an access token.
 */
func (i Config) GetProfile() (string, error) {

	err := i.validateGetProfileConfig()
	if err != nil {
		return "", fmt.Errorf("invalid parameters for profile: %v", err)
	}

	c, err := ims.NewClient(&ims.ClientConfig{
		URL: i.URL,
	})
	if err != nil {
		return "", fmt.Errorf("error creating the client: %v", err)
	}

	profile, err := c.GetProfile(&ims.GetProfileRequest{
		AccessToken: i.UserToken,
	})
	if err != nil {
		return "", err
	}

	return string(profile.Body), nil
}
