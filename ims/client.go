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
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Create a http client based on the received configuration.
func (i Config) httpClient() (*http.Client, error) {

	if i.ProxyURL != "" {
		p, err := url.Parse(i.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("proxy provided but its URL is malformed")
		}
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.Proxy = http.ProxyURL(p)
		if i.ProxyIgnoreTLS {
			t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		client := &http.Client{
			Timeout:   time.Duration(i.Timeout) * time.Second,
			Transport: t,
		}
		return client, nil
	}
	client := &http.Client{
		Timeout: time.Duration(i.Timeout) * time.Second,
	}
	return client, nil
}
