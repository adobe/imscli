// Copyright 2026 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

// Promise Token flows use token v4 and the "promise" / "promise_exchange" grant
// types per IMS Promise Token documentation. The promise-token grant exchanges an
// authenticating access token for a long-lived promise token, and the
// promise-exchange grant redeems a promise token for a fresh access token.

package ims

import (
	"fmt"

	"github.com/adobe/ims-go/ims"
)

// PromiseTokenInfo holds the response data from a promise-token exchange.
type PromiseTokenInfo struct {
	PromiseToken string
	TokenType    string
	Scope        string
	ExpiresIn    int
}

// PromiseExchangeInfo holds the response data from a promise-exchange.
type PromiseExchangeInfo struct {
	AccessToken           string
	PromiseToken          string
	PromiseTokenID        string
	TokenType             string
	Scope                 string
	ExpiresIn             int
	PromiseTokenExpiresIn int
}

func (i Config) validatePromiseTokenConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case !validateURL(i.URL):
		return fmt.Errorf("invalid IMS base URL parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client ID parameter")
	case i.ClientSecret == "":
		return fmt.Errorf("missing client secret parameter")
	case i.PromiseDefinitionID == "":
		return fmt.Errorf("missing promise definition ID parameter")
	case i.AccessToken == "":
		return fmt.Errorf("missing authenticating token parameter")
	default:
		return nil
	}
}

// PromiseTokenExchange exchanges an authenticating access token for a promise token.
func (i Config) PromiseTokenExchange() (PromiseTokenInfo, error) {
	if err := i.validatePromiseTokenConfig(); err != nil {
		return PromiseTokenInfo{}, fmt.Errorf("invalid parameters for promise token exchange: %w", err)
	}

	c, err := i.newIMSClient()
	if err != nil {
		return PromiseTokenInfo{}, fmt.Errorf("error creating the IMS client: %w", err)
	}

	r, err := c.PromiseToken(&ims.PromiseTokenRequest{
		ClientID:            i.ClientID,
		ClientSecret:        i.ClientSecret,
		PromiseDefinitionID: i.PromiseDefinitionID,
		AuthenticatingToken: i.AccessToken,
		Scopes:              i.Scopes,
	})
	if err != nil {
		return PromiseTokenInfo{}, fmt.Errorf("error during the promise token exchange: %w", err)
	}

	return PromiseTokenInfo{
		PromiseToken: r.PromiseToken,
		TokenType:    r.TokenType,
		Scope:        r.Scope,
		ExpiresIn:    int(r.ExpiresIn.Seconds()),
	}, nil
}

func (i Config) validatePromiseExchangeConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case !validateURL(i.URL):
		return fmt.Errorf("invalid IMS base URL parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client ID parameter")
	case i.ClientSecret == "":
		return fmt.Errorf("missing client secret parameter")
	case i.PromiseToken == "":
		return fmt.Errorf("missing promise token parameter")
	default:
		return nil
	}
}

// PromiseExchange redeems a promise token for a fresh access token.
func (i Config) PromiseExchange() (PromiseExchangeInfo, error) {
	if err := i.validatePromiseExchangeConfig(); err != nil {
		return PromiseExchangeInfo{}, fmt.Errorf("invalid parameters for promise exchange: %w", err)
	}

	c, err := i.newIMSClient()
	if err != nil {
		return PromiseExchangeInfo{}, fmt.Errorf("error creating the IMS client: %w", err)
	}

	r, err := c.PromiseExchange(&ims.PromiseExchangeRequest{
		ClientID:     i.ClientID,
		ClientSecret: i.ClientSecret,
		PromiseToken: i.PromiseToken,
		Scopes:       i.Scopes,
	})
	if err != nil {
		return PromiseExchangeInfo{}, fmt.Errorf("error during the promise exchange: %w", err)
	}

	return PromiseExchangeInfo{
		AccessToken:           r.AccessToken,
		PromiseToken:          r.PromiseToken,
		PromiseTokenID:        r.PromiseTokenID,
		TokenType:             r.TokenType,
		Scope:                 r.Scope,
		ExpiresIn:             int(r.ExpiresIn.Seconds()),
		PromiseTokenExpiresIn: int(r.PromiseTokenExpiresIn.Seconds()),
	}, nil
}
