// Copyright 2023 Adobe. All rights reserved.
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
)

/*
 * AuthorizeOAuth : Login for the OAuth server to server IMS flow
 */
func (i Config) AuthorizeOAuth() (string, error) {

	httpClient, err := i.httpClient()
	if err != nil {
		return "", fmt.Errorf("error creating the HTTP Client: %v", err)
	}

	c, err := ims.NewClient(&ims.ClientConfig{
		URL:    i.URL,
		Client: httpClient,
	})
	if err != nil {
		return "", fmt.Errorf("create client: %v", err)
	}

	r, err := c.Token(&ims.TokenRequest{
		ClientID:     i.ClientID,
		ClientSecret: i.ClientSecret,
		Scope:        i.Scopes,
		GrantType:    "client_credentials",
	})
	if err != nil {
		return "", fmt.Errorf("request token: %v", err)
	}

	return r.AccessToken, nil
}
