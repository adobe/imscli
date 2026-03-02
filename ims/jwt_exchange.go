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
	"os"
	"strings"
	"time"

	"github.com/adobe/ims-go/ims"
)

// jwtExpiration is the lifetime of JWT assertions used for the exchange flow.
const jwtExpiration = 30 * time.Minute

func (i Config) validateAuthorizeJWTExchangeConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case !validateURL(i.URL):
		return fmt.Errorf("invalid IMS base URL parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client ID parameter")
	case i.ClientSecret == "":
		return fmt.Errorf("missing client secret parameter")
	case i.PrivateKeyPath == "":
		return fmt.Errorf("missing private key path parameter")
	case i.Organization == "":
		return fmt.Errorf("missing organization parameter")
	case i.Account == "":
		return fmt.Errorf("missing account parameter")
	default:
		return nil
	}
}

// AuthorizeJWTExchange performs the JWT Bearer exchange flow.
func (i Config) AuthorizeJWTExchange() (TokenInfo, error) {

	if err := i.validateAuthorizeJWTExchangeConfig(); err != nil {
		return TokenInfo{}, fmt.Errorf("invalid parameters for JWT exchange: %w", err)
	}

	c, err := i.newIMSClient()
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error creating the IMS client: %w", err)
	}

	key, err := os.ReadFile(i.PrivateKeyPath)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error reading private key file %s: %w", i.PrivateKeyPath, err)
	}
	defer func() {
		for i := range key {
			key[i] = 0
		}
	}()

	// 	Metascopes are passed as generic claims with the format map[string]interface{}
	//  where the strings are in the form: baseIMSUrl/s/metascope
	//  and the interface{} is 'true'

	baseURL := strings.TrimRight(i.URL, "/")
	claims := make(map[string]interface{})
	for _, metascope := range i.Metascopes {
		claims[fmt.Sprintf("%s/s/%s", baseURL, metascope)] = true
	}

	r, err := c.ExchangeJWT(&ims.ExchangeJWTRequest{
		PrivateKey:   key,
		Expiration:   time.Now().Add(jwtExpiration),
		Issuer:       i.Organization,
		Subject:      i.Account,
		ClientID:     i.ClientID,
		ClientSecret: i.ClientSecret,
		Claims:       claims,
	})
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error exchanging JWT: %w", err)
	}

	return TokenInfo{
		AccessToken: r.AccessToken,
	}, nil
}
