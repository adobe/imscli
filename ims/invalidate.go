// Copyright 2021 Adobe. All rights reserved.
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

/*
 * Invalidate a token
 */
func (i Config) validateInvalidateTokenConfig() error {

	switch {
	case i.ClientID == "":
		return fmt.Errorf("missing clientID parameter")
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case i.AccessToken != "":
		log.Println("access token will be invalidated")
		return nil
	case i.RefreshToken != "":
		log.Println("refresh token will be invalidated")
		return nil
	case i.DeviceToken != "":
		log.Println("device token will be invalidated")
		return nil
	case i.ServiceToken != "":
		log.Println("authorization code will be invalidated")
		if i.ClientSecret == "" {
			return fmt.Errorf("missing client secret, mandatory to invalidate service token")
		}
		return nil
	default:
		return fmt.Errorf("no token has been found for invalidation")
	}
}

// InvalidateToken Invalidates the token provided in the configuration using the IMS API.
func (i Config) InvalidateToken() error {
	// Perform parameter validation
	err := i.validateValidateTokenConfig()
	if err != nil {
		return fmt.Errorf("incomplete parameters for token invalidation: %v", err)
	}

	httpClient, err := i.httpClient()
	if err != nil {
		return fmt.Errorf("error creating the HTTP Client: %v", err)
	}

	c, err := ims.NewClient(&ims.ClientConfig{
		URL:    i.URL,
		Client: httpClient,
	})
	if err != nil {
		return fmt.Errorf("create client: %v", err)
	}

	var token string
	var tokenType ims.TokenType

	switch {
	case i.AccessToken != "":
		token = i.AccessToken
		tokenType = ims.AccessToken
	case i.RefreshToken != "":
		token = i.RefreshToken
		tokenType = ims.RefreshToken
	case i.DeviceToken != "":
		token = i.DeviceToken
		tokenType = ims.DeviceToken
	case i.ServiceToken != "":
		token = i.ServiceToken
		tokenType = ims.ServiceToken
	default:
		return fmt.Errorf("unexpected error, broken parameter validation")
	}

	err = c.InvalidateToken(&ims.InvalidateTokenRequest{
		Token:        token,
		Type:         tokenType,
		ClientID:     i.ClientID,
		Cascading:    i.Cascading,
		ClientSecret: i.ClientSecret,
	})
	if err != nil {
		return fmt.Errorf("error during token invalidation: %v", err)
	}

	return nil
}
