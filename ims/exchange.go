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
	"time"

	"github.com/adobe/ims-go/ims"
)

func (i Config) validateClusterExchangeConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client ID parameter")
	case i.ClientSecret == "":
		return fmt.Errorf("missing client secret parameter")
	case i.AccessToken == "":
		return fmt.Errorf("missing access token parameter")
	case i.UserID != "" && i.Organization != "":
		return fmt.Errorf("can't perform the request with user ID and IMS Organization at the same time")
	case i.UserID == "" && i.Organization == "":
		return fmt.Errorf("missing user ID or IMS Organization parameter")
	default:
		return nil
	}
}

// ClusterExchange performs the Cluster Access Token Exchange grant flow
func (i Config) ClusterExchange() (TokenInfo, error) {

	if err := i.validateClusterExchangeConfig(); err != nil {
		return TokenInfo{}, fmt.Errorf("invalid parameters for cluster exchange: %w", err)
	}

	httpClient, err := i.httpClient()
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error creating the HTTP Client: %w", err)
	}

	c, err := ims.NewClient(&ims.ClientConfig{
		URL:    i.URL,
		Client: httpClient,
	})
	if err != nil {
		return TokenInfo{}, fmt.Errorf("create client: %w", err)
	}

	r, err := c.ClusterExchange(&ims.ClusterExchangeRequest{
		ClientID:     i.ClientID,
		ClientSecret: i.ClientSecret,
		UserToken:    i.AccessToken,
		UserID:       i.UserID,
		OrgID:        i.Organization,
		Scopes:       i.Scopes,
	})
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error during the cluster exchange: %w", err)
	}

	return TokenInfo{
		AccessToken: r.AccessToken,
		Expires:     int(r.ExpiresIn * time.Millisecond),
	}, nil
}
