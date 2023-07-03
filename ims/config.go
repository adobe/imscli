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
	"net/url"
)

type Config struct {
	URL               string
	ClientID          string
	ClientSecret      string
	ServiceToken      string
	PrivateKeyPath    string
	Organization      string
	Account           string
	Scopes            []string
	Metascopes        []string
	AccessToken       string
	RefreshToken      string
	DeviceToken       string
	AuthorizationCode string
	ProfileApiVersion string
	OrgsApiVersion    string
	Timeout           int
	ProxyURL          string
	ProxyIgnoreTLS    bool
	PublicClient      bool
	UserID            string
	Cascading         bool
	Token             string
	Port              int
}

type TokenInfo struct {
	AccessToken string
	Expires     int //(response.ExpiresIn * time.Millisecond),
	Valid       bool
	Info        string
}

func validateURL(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}

	switch {
	case parsedURL.Scheme == "":
		return false
	case parsedURL.Host == "":
		return false
	default:
		return true
	}
}
