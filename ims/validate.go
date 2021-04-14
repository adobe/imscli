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
 * Validate that the config includes:
 * - One clientID
 * - One token to validate
 */

func (i Config) validateValidateTokenConfig() error {

	switch {
	case i.ClientID == "":
		return fmt.Errorf("missing clientID parameter")
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case i.AccessToken != "":
		log.Println("access token will be validated")
		return nil
	case i.RefreshToken != "":
		log.Println("refresh token will be validated")
		return nil
	case i.DeviceToken != "":
		log.Println("device token will be validated")
		return nil
	case i.AuthorizationCode != "":
		log.Println("authorization code will be validated")
		return nil
	default:
		return fmt.Errorf("no token type has been found for validation")
	}
}

// ValidateToken Validates the token provided in the configuration using the IMS API.
// Return the endpoint response or an error.
func (i Config) ValidateToken() (TokenInfo, error) {
	// Perform parameter validation
	err := i.validateValidateTokenConfig()
	if err != nil {
		return TokenInfo{}, fmt.Errorf("invalid parameters for token validation: %v", err)
	}

	httpClient, err := i.httpClient()
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error creating the HTTP Client: %v", err)
	}

	c, err := ims.NewClient(&ims.ClientConfig{
		URL:    i.URL,
		Client: httpClient,
	})
	if err != nil {
		return TokenInfo{}, fmt.Errorf("create client: %v", err)
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
	case i.AuthorizationCode != "":
		token = i.AuthorizationCode
		tokenType = ims.AuthorizationCode
	default:
		return TokenInfo{}, fmt.Errorf("unexpected error, broken validation")
	}

	r, err := c.ValidateToken(&ims.ValidateTokenRequest{
		Token:    token,
		Type:     tokenType,
		ClientID: i.ClientID,
	})
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error during token validation: %v", err)
	}

	return TokenInfo{
		Valid: r.Valid,
		Info:  string(r.Body),
	}, nil
}
