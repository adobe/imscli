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
	"github.com/adobe/ims-go/ims"
	"time"
)

func (i Config) validateRefreshConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client ID parameter")
	case i.ClientSecret == "":
		return fmt.Errorf("missing client secret parameter")
	case i.RefreshToken == "":
		return fmt.Errorf("missing refresh token parameter")
	default:
		return nil
	}
}

// Refresh performs the refresh token flow.
func (i Config) Refresh() (RefreshInfo, error) {

	if err := i.validateRefreshConfig(); err != nil {
		return RefreshInfo{}, fmt.Errorf("invalid parameters for token refresh: %v", err)
	}

	httpClient, err := i.httpClient()
	if err != nil {
		return RefreshInfo{}, fmt.Errorf("error creating the HTTP Client: %v", err)
	}

	c, err := ims.NewClient(&ims.ClientConfig{
		URL:    i.URL,
		Client: httpClient,
	})
	if err != nil {
		return RefreshInfo{}, fmt.Errorf("create client: %v", err)
	}

	r, err := c.RefreshToken(&ims.RefreshTokenRequest{
		ClientID:     i.ClientID,
		ClientSecret: i.ClientSecret,
		RefreshToken: i.RefreshToken,
		Scope:        i.Scopes,
	})
	if err != nil {
		return RefreshInfo{}, fmt.Errorf("error during the token refresh: %v", err)
	}

	return RefreshInfo{
		TokenInfo: TokenInfo{
			AccessToken: r.AccessToken,
			Expires:     int(r.ExpiresIn * time.Second),
		},
		RefreshToken: r.RefreshToken,
	}, nil
}
