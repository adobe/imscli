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

func (i Config) AuthorizeJWTExchange() (TokenInfo, error) {

	c, err := i.newIMSClient()
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error creating the IMS client: %w", err)
	}

	key, err := os.ReadFile(i.PrivateKeyPath)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error read private key file: %s, %w", i.PrivateKeyPath, err)
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
		Expiration:   time.Now().Add(time.Minute * 30),
		Issuer:       i.Organization,
		Subject:      i.Account,
		ClientID:     i.ClientID,
		ClientSecret: i.ClientSecret,
		Claims:       claims,
	})
	if err != nil {
		return TokenInfo{}, fmt.Errorf("exchange JWT: %w", err)
	}

	return TokenInfo{
		AccessToken: r.AccessToken,
	}, nil
}
