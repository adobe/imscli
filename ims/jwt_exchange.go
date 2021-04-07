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
	"io/ioutil"
	"time"

	"github.com/adobe/ims-go/ims"
)

func (i Config) AuthorizeJWTExchange() (TokenInfo, error) {

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

	key, err := ioutil.ReadFile(i.PrivateKeyPath)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("read private key file: %v", err)
	}

	// 	Metascopes are passed as generic claims with the format map[string]interface{}
	//  where the strings are in the form: baseIMSUrl/s/metascope
	//  and the interface{} is 'true'

	claims := make(map[string]interface{})
	for _, metascope := range i.Metascopes {
		claims[fmt.Sprintf("%s/s/%s", i.URL, metascope)] = true
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
		return TokenInfo{}, fmt.Errorf("exchange JWT: %v", err)
	}

	return TokenInfo{
		AccessToken: r.AccessToken,
		Expires:     int(r.ExpiresIn * time.Millisecond),
	}, nil
}
