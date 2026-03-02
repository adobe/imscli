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
)

func (i Config) validateAuthorizeServiceConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case !validateURL(i.URL):
		return fmt.Errorf("invalid IMS base URL parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client ID parameter")
	case i.ClientSecret == "":
		return fmt.Errorf("missing client secret parameter")
	case i.AuthorizationCode == "":
		return fmt.Errorf("missing authorization code parameter")
	default:
		return nil
	}
}

// AuthorizeService performs the service-to-service IMS authorization flow.
func (i Config) AuthorizeService() (string, error) {

	if err := i.validateAuthorizeServiceConfig(); err != nil {
		return "", fmt.Errorf("invalid parameters for service authorization: %w", err)
	}

	c, err := i.newIMSClient()
	if err != nil {
		return "", fmt.Errorf("error creating the IMS client: %w", err)
	}

	r, err := c.Token(&ims.TokenRequest{
		ClientID:     i.ClientID,
		ClientSecret: i.ClientSecret,
		Code:         i.AuthorizationCode,
	})
	if err != nil {
		return "", fmt.Errorf("error requesting token: %w", err)
	}

	return r.AccessToken, nil
}
