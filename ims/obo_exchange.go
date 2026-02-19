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

// OBO (On-Behalf-Of) token exchange security constraints:
// - Do NOT send OBO access tokens to frontend clients; they are for backend-to-backend use only.
// - Short TTL: OBO tokens typically expire in 5 minutes by default.
// - Subject token restrictions: only user access tokens are accepted; ServiceTokens and
//   impersonation tokens must not be used as the subject token.
// - Scope boundary: requested scopes cannot exceed the client's configured scopes.
// - Audit trail: the full actor chain is preserved in the act claim of the issued token.

// grant type for OBO at IMS token endpoint (adjust if Adobe IMS uses a different value)
const oboGrantType = "on_behalf_of"

func (i Config) validateOBOExchangeConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client ID parameter")
	case i.ClientSecret == "":
		return fmt.Errorf("missing client secret parameter")
	case i.AccessToken == "":
		return fmt.Errorf("missing access token parameter (user token only; do not use service or impersonation tokens)")
	case i.ServiceToken != "":
		return fmt.Errorf("OBO exchange requires a user access token; do not use service token as subject token")
	default:
		return nil
	}
}

// OBOExchange performs the On-Behalf-Of token exchange.
// It exchanges a user's access token for a new token that can be used by the backend to call
// downstream APIs on behalf of that user. The returned token must only be used server-side.
func (i Config) OBOExchange() (TokenInfo, error) {
	if err := i.validateOBOExchangeConfig(); err != nil {
		return TokenInfo{}, fmt.Errorf("invalid parameters for OBO exchange: %v", err)
	}

	httpClient, err := i.httpClient()
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error creating the HTTP Client: %v", err)
	}

	data := url.Values{}
	data.Set("grant_type", oboGrantType)
	data.Set("client_secret", i.ClientSecret)
	data.Set("subject_token", i.AccessToken)
	data.Set("subject_token_type", "urn:ietf:params:oauth:token-type:access_token")
	if len(i.Scopes) > 0 && (len(i.Scopes) != 1 || i.Scopes[0] != "") {
		data.Set("scope", strings.Join(i.Scopes, ","))
	}

	tokenURL := fmt.Sprintf("%s/ims/token/v3?client_id=%s", strings.TrimSuffix(i.URL, "/"), url.QueryEscape(i.ClientID))
	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error creating OBO request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error during OBO exchange: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error reading OBO response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return TokenInfo{}, fmt.Errorf("OBO exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return TokenInfo{}, fmt.Errorf("error decoding OBO response: %v", err)
	}

	expiresMs := 0
	if out.ExpiresIn > 0 {
		expiresMs = out.ExpiresIn * int(time.Second/time.Millisecond)
	}

	return TokenInfo{
		AccessToken: out.AccessToken,
		Expires:     expiresMs,
	}, nil
}
