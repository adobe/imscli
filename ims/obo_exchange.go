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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// - Subject token restrictions: only user access tokens are accepted; ServiceTokens and
//   impersonation tokens must not be used as the subject token.
// - Scope boundary: requested scopes cannot exceed the client's configured scopes.
// - Audit trail: the full actor chain is preserved in the act claim of the issued token.

// OBO uses token v4 and RFC 8693 grant type per IMS OBO documentation.
const defaultOBOGrantType = "urn:ietf:params:oauth:grant-type:token-exchange"

func (i Config) validateOBOExchangeConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client ID parameter")
	case i.ClientSecret == "":
		return fmt.Errorf("missing client secret parameter")
	case i.AccessToken == "":
		return fmt.Errorf("missing access token parameter (only access tokens are accepted)")
	case len(i.Scopes) == 0 || (len(i.Scopes) == 1 && i.Scopes[0] == ""):
		return fmt.Errorf("scopes are required for On-Behalf-Of exchange")
	default:
		return nil
	}
}

func (i Config) OBOExchange() (TokenInfo, error) {
	if err := i.validateOBOExchangeConfig(); err != nil {
		return TokenInfo{}, fmt.Errorf("invalid parameters for On-Behalf-Of exchange: %v", err)
	}

	httpClient, err := i.httpClient()
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error creating the HTTP Client: %v", err)
	}

	data := url.Values{}
	data.Set("grant_type", defaultOBOGrantType)
	data.Set("client_id", i.ClientID)
	data.Set("client_secret", i.ClientSecret)
	data.Set("subject_token", i.AccessToken)
	data.Set("subject_token_type", "urn:ietf:params:oauth:token-type:access_token")
	data.Set("requested_token_type", "urn:ietf:params:oauth:token-type:access_token")
	data.Set("scope", strings.Join(i.Scopes, ","))

	// OBO Token Exchange requires /ims/token/v4 (v3 does not support this grant type).
	tokenURL := fmt.Sprintf("%s/ims/token/v4?client_id=%s", strings.TrimSuffix(i.URL, "/"), url.QueryEscape(i.ClientID))
	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error creating On-Behalf-Of request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error during On-Behalf-Of exchange: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error reading On-Behalf-Of response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("On-Behalf-Of exchange failed (status %d): %s", resp.StatusCode, string(body))
		if resp.StatusCode == http.StatusBadRequest {
			if strings.Contains(string(body), "invalid_scope") {
				errMsg += " — IMS may be rejecting the subject token's scopes for this client. Ensure the client has Token exchange enabled and allowed scopes in the portal, or try a user token obtained with fewer scopes."
			}
		}
		return TokenInfo{}, fmt.Errorf("%s", errMsg)
	}

	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return TokenInfo{}, fmt.Errorf("error decoding On-Behalf-Of response: %v", err)
	}

	expiresMs := out.ExpiresIn * int(time.Second/time.Millisecond)

	return TokenInfo{
		AccessToken: out.AccessToken,
		Expires:     expiresMs,
	}, nil
}
