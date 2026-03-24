// Copyright 2021 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

// - Subject token restrictions: only user access tokens are accepted; ServiceTokens and
//   impersonation tokens must not be used as the subject token.
// - Scope boundary: requested scopes cannot exceed the client's configured scopes.
// - Audit trail: the full actor chain is preserved in the act claim of the issued token.

// OBO uses token v4 and RFC 8693 grant type per IMS OBO documentation.

package ims

import (
	"fmt"

	"github.com/adobe/ims-go/ims"
)

func (i Config) validateOBOExchangeConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case !validateURL(i.URL):
		return fmt.Errorf("invalid IMS base URL parameter")
	case i.AccessToken == "":
		return fmt.Errorf("missing access token parameter")
	default:
		return nil
	}
}

func (i Config) OBOExchange() (TokenInfo, error) {

	if err := i.validateOBOExchangeConfig(); err != nil {
		return TokenInfo{}, fmt.Errorf("invalid parameters for On-Behalf-Of exchange: %w", err)
	}

	c, err := i.newIMSClient()
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error creating the IMS client: %w", err)
	}

	r, err := c.OBOExchange(&ims.OBOExchangeRequest{
		ClientID:     i.ClientID,
		ClientSecret: i.ClientSecret,
		SubjectToken: i.AccessToken,
		Scopes:       i.Scopes,
	})
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error during the On-Behalf-Of exchange: %w", err)
	}

	return TokenInfo{
		AccessToken: r.AccessToken,
	}, nil
}
