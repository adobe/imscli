// Copyright 2026 Adobe. All rights reserved.
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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateAuthorizeImplicitConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{
			name:    "valid",
			config:  Config{URL: "https://ims.example.com", ClientID: "c", Scopes: []string{"openid"}, Port: 8888, RedirectURI: DefaultImplicitRedirectURI},
			wantErr: "",
		},
		{
			name:    "missing URL",
			config:  Config{ClientID: "c", Scopes: []string{"openid"}, Port: 8888, RedirectURI: DefaultImplicitRedirectURI},
			wantErr: "missing IMS base URL",
		},
		{
			name:    "invalid URL",
			config:  Config{URL: "://bad", ClientID: "c", Scopes: []string{"openid"}, Port: 8888, RedirectURI: DefaultImplicitRedirectURI},
			wantErr: "unable to parse URL",
		},
		{
			name:    "missing scopes",
			config:  Config{URL: "https://ims.example.com", ClientID: "c", Port: 8888, RedirectURI: DefaultImplicitRedirectURI},
			wantErr: "missing scopes",
		},
		{
			name:    "empty scope",
			config:  Config{URL: "https://ims.example.com", ClientID: "c", Scopes: []string{""}, Port: 8888, RedirectURI: DefaultImplicitRedirectURI},
			wantErr: "missing scopes",
		},
		{
			name:    "missing clientID",
			config:  Config{URL: "https://ims.example.com", Scopes: []string{"openid"}, Port: 8888, RedirectURI: DefaultImplicitRedirectURI},
			wantErr: "missing client id",
		},
		{
			name:    "missing port",
			config:  Config{URL: "https://ims.example.com", ClientID: "c", Scopes: []string{"openid"}, RedirectURI: DefaultImplicitRedirectURI},
			wantErr: "missing or invalid port",
		},
		{
			name:    "negative port",
			config:  Config{URL: "https://ims.example.com", ClientID: "c", Scopes: []string{"openid"}, Port: -1, RedirectURI: DefaultImplicitRedirectURI},
			wantErr: "missing or invalid port",
		},
		{
			name:    "missing redirect URI",
			config:  Config{URL: "https://ims.example.com", ClientID: "c", Scopes: []string{"openid"}, Port: 8888},
			wantErr: "missing redirect URI",
		},
		{
			name:    "invalid redirect URI",
			config:  Config{URL: "https://ims.example.com", ClientID: "c", Scopes: []string{"openid"}, Port: 8888, RedirectURI: "://bad"},
			wantErr: "unable to parse redirect URI",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateAuthorizeImplicitConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestImplicitCapture(t *testing.T) {
	const state = "expected-state-abc"

	tests := []struct {
		name      string
		url       string
		wantToken string
		wantErr   string
	}{
		{
			name:      "happy path",
			url:       "/?access_token=tok&state=" + state + "&expires_in=3600&token_type=bearer",
			wantToken: "tok",
		},
		{
			name:    "state mismatch",
			url:     "/?access_token=tok&state=attacker",
			wantErr: "state mismatch",
		},
		{
			name:    "missing state",
			url:     "/?access_token=tok",
			wantErr: "state mismatch",
		},
		{
			name:    "IMS error param",
			url:     "/?error=access_denied&error_description=user+rejected&state=" + state,
			wantErr: "authorization error: access_denied: user rejected",
		},
		{
			name:    "missing token",
			url:     "/?state=" + state,
			wantErr: "missing access_token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resCh := make(chan *TokenInfo, 1)
			errCh := make(chan error, 1)
			h := &implicitHandler{
				expectedState: state,
				resCh:         resCh,
				errCh:         errCh,
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			h.capture(rec, req)

			select {
			case res := <-resCh:
				if tt.wantErr != "" {
					t.Fatalf("expected error %q, got token %q", tt.wantErr, res.AccessToken)
				}
				if res.AccessToken != tt.wantToken {
					t.Errorf("token: got %q, want %q", res.AccessToken, tt.wantToken)
				}
			case err := <-errCh:
				if tt.wantErr == "" {
					t.Fatalf("expected token, got error %v", err)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error: got %q, want it to contain %q", err.Error(), tt.wantErr)
				}
			default:
				t.Fatalf("handler neither sent token nor error")
			}

			if got := rec.Result().Header.Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
				t.Errorf("Content-Type: got %q, want text/html prefix", got)
			}
		})
	}
}

func TestRandomStateProducesDistinctValues(t *testing.T) {
	a, err := randomState()
	if err != nil {
		t.Fatalf("randomState: %v", err)
	}
	b, err := randomState()
	if err != nil {
		t.Fatalf("randomState: %v", err)
	}
	if a == b {
		t.Errorf("randomState produced identical values across calls")
	}
	if len(a) < 40 {
		t.Errorf("randomState too short: %d chars", len(a))
	}
}
