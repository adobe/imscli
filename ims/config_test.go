// Copyright 2025 Adobe. All rights reserved.
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
	"math/rand"
	"testing"
	"time"
)

func TestValidateDecodeTokenConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid", config: Config{Token: "a.b.c"}, wantErr: ""},
		{name: "missing token", config: Config{}, wantErr: "missing token parameter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateDecodeTokenConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateValidateTokenConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid with access token", config: Config{ClientID: "c", URL: "https://ims.example.com", AccessToken: "tok"}, wantErr: ""},
		{name: "valid with refresh token", config: Config{ClientID: "c", URL: "https://ims.example.com", RefreshToken: "tok"}, wantErr: ""},
		{name: "valid with device token", config: Config{ClientID: "c", URL: "https://ims.example.com", DeviceToken: "tok"}, wantErr: ""},
		{name: "valid with authorization code", config: Config{ClientID: "c", URL: "https://ims.example.com", AuthorizationCode: "tok"}, wantErr: ""},
		{name: "missing clientID", config: Config{URL: "https://ims.example.com", AccessToken: "tok"}, wantErr: "missing clientID"},
		{name: "missing URL", config: Config{ClientID: "c", AccessToken: "tok"}, wantErr: "missing IMS base URL"},
		{name: "missing token", config: Config{ClientID: "c", URL: "https://ims.example.com"}, wantErr: "no token type has been found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateValidateTokenConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateInvalidateTokenConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid with access token", config: Config{ClientID: "c", URL: "https://ims.example.com", AccessToken: "tok"}, wantErr: ""},
		{name: "valid with refresh token", config: Config{ClientID: "c", URL: "https://ims.example.com", RefreshToken: "tok"}, wantErr: ""},
		{name: "valid with device token", config: Config{ClientID: "c", URL: "https://ims.example.com", DeviceToken: "tok"}, wantErr: ""},
		{name: "valid with service token", config: Config{ClientID: "c", URL: "https://ims.example.com", ServiceToken: "tok", ClientSecret: "s"}, wantErr: ""},
		{name: "service token without secret", config: Config{ClientID: "c", URL: "https://ims.example.com", ServiceToken: "tok"}, wantErr: "missing client secret"},
		{name: "missing clientID", config: Config{URL: "https://ims.example.com", AccessToken: "tok"}, wantErr: "missing clientID"},
		{name: "missing URL", config: Config{ClientID: "c", AccessToken: "tok"}, wantErr: "missing IMS base URL"},
		{name: "missing token", config: Config{ClientID: "c", URL: "https://ims.example.com"}, wantErr: "no token has been found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateInvalidateTokenConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateGetProfileConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid", config: Config{ProfileAPIVersion: "v2", AccessToken: "tok", URL: "https://ims.example.com"}, wantErr: ""},
		{name: "invalid version", config: Config{ProfileAPIVersion: "v99", AccessToken: "tok", URL: "https://ims.example.com"}, wantErr: "invalid API version"},
		{name: "missing access token", config: Config{ProfileAPIVersion: "v1", URL: "https://ims.example.com"}, wantErr: "missing access token"},
		{name: "missing URL", config: Config{ProfileAPIVersion: "v1", AccessToken: "tok"}, wantErr: "missing IMS base URL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateGetProfileConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateGetAdminProfileConfig(t *testing.T) {
	validConfig := Config{
		ProfileAPIVersion: "v2",
		ServiceToken:      "tok",
		URL:               "https://ims.example.com",
		ClientID:          "c",
		Guid:              "g",
		AuthSrc:           "a",
	}
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid", config: validConfig, wantErr: ""},
		{name: "invalid version", config: withField(validConfig, func(c *Config) { c.ProfileAPIVersion = "v99" }), wantErr: "invalid API version"},
		{name: "missing service token", config: withField(validConfig, func(c *Config) { c.ServiceToken = "" }), wantErr: "missing service token"},
		{name: "missing URL", config: withField(validConfig, func(c *Config) { c.URL = "" }), wantErr: "missing IMS base URL"},
		{name: "missing clientID", config: withField(validConfig, func(c *Config) { c.ClientID = "" }), wantErr: "missing client ID"},
		{name: "missing guid", config: withField(validConfig, func(c *Config) { c.Guid = "" }), wantErr: "missing guid"},
		{name: "missing authSrc", config: withField(validConfig, func(c *Config) { c.AuthSrc = "" }), wantErr: "missing auth source"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateGetAdminProfileConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateGetOrganizationsConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid", config: Config{OrgsAPIVersion: "v5", AccessToken: "tok", URL: "https://ims.example.com"}, wantErr: ""},
		{name: "invalid version", config: Config{OrgsAPIVersion: "v99", AccessToken: "tok", URL: "https://ims.example.com"}, wantErr: "invalid API version"},
		{name: "missing access token", config: Config{OrgsAPIVersion: "v5", URL: "https://ims.example.com"}, wantErr: "missing access token"},
		{name: "missing URL", config: Config{OrgsAPIVersion: "v5", AccessToken: "tok"}, wantErr: "missing IMS base URL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateGetOrganizationsConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateGetAdminOrganizationsConfig(t *testing.T) {
	validConfig := Config{
		OrgsAPIVersion: "v5",
		ServiceToken:   "tok",
		URL:            "https://ims.example.com",
		ClientID:       "c",
		Guid:           "g",
		AuthSrc:        "a",
	}
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid", config: validConfig, wantErr: ""},
		{name: "invalid version", config: withField(validConfig, func(c *Config) { c.OrgsAPIVersion = "v99" }), wantErr: "invalid API version"},
		{name: "missing service token", config: withField(validConfig, func(c *Config) { c.ServiceToken = "" }), wantErr: "missing service token"},
		{name: "missing URL", config: withField(validConfig, func(c *Config) { c.URL = "" }), wantErr: "missing IMS base URL"},
		{name: "missing clientID", config: withField(validConfig, func(c *Config) { c.ClientID = "" }), wantErr: "missing client ID"},
		{name: "missing guid", config: withField(validConfig, func(c *Config) { c.Guid = "" }), wantErr: "missing guid"},
		{name: "missing authSrc", config: withField(validConfig, func(c *Config) { c.AuthSrc = "" }), wantErr: "missing auth source"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateGetAdminOrganizationsConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateClusterExchangeConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid with org", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s", AccessToken: "tok", Organization: "org"}, wantErr: ""},
		{name: "valid with userID", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s", AccessToken: "tok", UserID: "u"}, wantErr: ""},
		{name: "missing URL", config: Config{ClientID: "c", ClientSecret: "s", AccessToken: "tok", Organization: "org"}, wantErr: "missing IMS base URL"},
		{name: "missing clientID", config: Config{URL: "https://ims.example.com", ClientSecret: "s", AccessToken: "tok", Organization: "org"}, wantErr: "missing client ID"},
		{name: "missing secret", config: Config{URL: "https://ims.example.com", ClientID: "c", AccessToken: "tok", Organization: "org"}, wantErr: "missing client secret"},
		{name: "missing access token", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s", Organization: "org"}, wantErr: "missing access token"},
		{name: "both org and userID", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s", AccessToken: "tok", Organization: "org", UserID: "u"}, wantErr: "can't perform the request"},
		{name: "neither org nor userID", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s", AccessToken: "tok"}, wantErr: "missing user ID or IMS Organization"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateClusterExchangeConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateRefreshConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s", RefreshToken: "tok"}, wantErr: ""},
		{name: "missing URL", config: Config{ClientID: "c", ClientSecret: "s", RefreshToken: "tok"}, wantErr: "missing IMS base URL"},
		{name: "missing clientID", config: Config{URL: "https://ims.example.com", ClientSecret: "s", RefreshToken: "tok"}, wantErr: "missing client ID"},
		{name: "missing secret", config: Config{URL: "https://ims.example.com", ClientID: "c", RefreshToken: "tok"}, wantErr: "missing client secret"},
		{name: "missing refresh token", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s"}, wantErr: "missing refresh token"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateRefreshConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateAuthorizeUserConfig(t *testing.T) {
	validConfig := Config{
		URL:          "https://ims.example.com",
		ClientID:     "c",
		ClientSecret: "s",
		Organization: "org",
		Scopes:       []string{"openid"},
		Port:         8888,
	}
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid", config: validConfig, wantErr: ""},
		{name: "valid public client", config: withField(validConfig, func(c *Config) { c.ClientSecret = ""; c.PublicClient = true }), wantErr: ""},
		{name: "missing URL", config: withField(validConfig, func(c *Config) { c.URL = "" }), wantErr: "missing IMS base URL"},
		{name: "invalid URL", config: withField(validConfig, func(c *Config) { c.URL = "not-a-url" }), wantErr: "unable to parse URL"},
		{name: "missing scopes", config: withField(validConfig, func(c *Config) { c.Scopes = nil }), wantErr: "missing scopes"},
		{name: "empty scopes", config: withField(validConfig, func(c *Config) { c.Scopes = []string{""} }), wantErr: "missing scopes"},
		{name: "missing clientID", config: withField(validConfig, func(c *Config) { c.ClientID = "" }), wantErr: "missing client id"},
		{name: "missing organization", config: withField(validConfig, func(c *Config) { c.Organization = "" }), wantErr: "missing organization"},
		{name: "missing secret non-public", config: withField(validConfig, func(c *Config) { c.ClientSecret = "" }), wantErr: "missing client secret"},
		{name: "missing port", config: withField(validConfig, func(c *Config) { c.Port = 0 }), wantErr: "missing or invalid port"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateAuthorizeUserConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateAuthorizeServiceConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s", AuthorizationCode: "code"}, wantErr: ""},
		{name: "missing URL", config: Config{ClientID: "c", ClientSecret: "s", AuthorizationCode: "code"}, wantErr: "missing IMS base URL"},
		{name: "missing clientID", config: Config{URL: "https://ims.example.com", ClientSecret: "s", AuthorizationCode: "code"}, wantErr: "missing client ID"},
		{name: "missing secret", config: Config{URL: "https://ims.example.com", ClientID: "c", AuthorizationCode: "code"}, wantErr: "missing client secret"},
		{name: "missing code", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s"}, wantErr: "missing authorization code"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateAuthorizeServiceConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateAuthorizeClientCredentialsConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s", Scopes: []string{"openid"}}, wantErr: ""},
		{name: "missing URL", config: Config{ClientID: "c", ClientSecret: "s", Scopes: []string{"openid"}}, wantErr: "missing IMS base URL"},
		{name: "missing clientID", config: Config{URL: "https://ims.example.com", ClientSecret: "s", Scopes: []string{"openid"}}, wantErr: "missing client ID"},
		{name: "missing secret", config: Config{URL: "https://ims.example.com", ClientID: "c", Scopes: []string{"openid"}}, wantErr: "missing client secret"},
		{name: "missing scopes", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s"}, wantErr: "missing scopes"},
		{name: "empty scopes", config: Config{URL: "https://ims.example.com", ClientID: "c", ClientSecret: "s", Scopes: []string{""}}, wantErr: "missing scopes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateAuthorizeClientCredentialsConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateRegisterConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{name: "valid", config: Config{URL: "https://ims.example.com", ClientName: "app", RedirectURIs: []string{"https://example.com/cb"}}, wantErr: ""},
		{name: "missing URL", config: Config{ClientName: "app", RedirectURIs: []string{"https://example.com/cb"}}, wantErr: "missing IMS base URL"},
		{name: "missing client name", config: Config{URL: "https://ims.example.com", RedirectURIs: []string{"https://example.com/cb"}}, wantErr: "missing client name"},
		{name: "missing redirect URIs", config: Config{URL: "https://ims.example.com", ClientName: "app"}, wantErr: "missing redirect URIs"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateRegisterConfig()
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestResolveToken(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantToken string
		wantErr   string
	}{
		{name: "access token", config: Config{AccessToken: "a"}, wantToken: "a"},
		{name: "refresh token", config: Config{RefreshToken: "r"}, wantToken: "r"},
		{name: "device token", config: Config{DeviceToken: "d"}, wantToken: "d"},
		{name: "service token", config: Config{ServiceToken: "s"}, wantToken: "s"},
		{name: "authorization code", config: Config{AuthorizationCode: "c"}, wantToken: "c"},
		{name: "no token", config: Config{}, wantErr: "no token provided"},
		{name: "multiple tokens", config: Config{AccessToken: "a", RefreshToken: "r"}, wantErr: "multiple tokens provided"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, _, err := tt.config.resolveToken()
			assertError(t, err, tt.wantErr)
			if err == nil && token != tt.wantToken {
				t.Errorf("resolveToken() = %q, want %q", token, tt.wantToken)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{name: "valid https", url: "https://ims.example.com", want: true},
		{name: "valid http", url: "http://localhost:8080", want: true},
		{name: "missing scheme", url: "ims.example.com", want: false},
		{name: "missing host", url: "https://", want: false},
		{name: "empty", url: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateURL(tt.url); got != tt.want {
				t.Errorf("validateURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

// withField returns a copy of the config with the given mutation applied.
func withField(c Config, mutate func(*Config)) Config {
	mutate(&c)
	return c
}

// assertError checks that err matches the expected substring, or is nil if wantErr is empty.
func assertError(t *testing.T, err error, wantErr string) {
	t.Helper()
	if wantErr == "" {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		return
	}
	if err == nil {
		t.Errorf("expected error containing %q, got nil", wantErr)
		return
	}
	if got := err.Error(); !contains(got, wantErr) {
		t.Errorf("error %q does not contain %q", got, wantErr)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// randomString generates a string of the given length with arbitrary bytes.
func randomString(rng *rand.Rand, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(rng.Intn(256))
	}
	return string(b)
}

// TestFuzzValidateURL generates random inputs for 10 seconds to verify that
// validateURL never panics regardless of input. Runs in parallel with other tests.
//
// For deeper exploration, use Go's built-in fuzz engine:
//
//	go test -fuzz=FuzzValidateURL -fuzztime=60s ./ims/
func TestFuzzValidateURL(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	deadline := time.After(10 * time.Second)
	iterations := 0

	for {
		select {
		case <-deadline:
			t.Logf("fuzz: %d iterations without panic", iterations)
			return
		default:
			input := randomString(rng, rng.Intn(512))
			_ = validateURL(input)
			iterations++
		}
	}
}

// FuzzValidateURL is a standard Go fuzz target for deeper exploration.
// Run manually: go test -fuzz=FuzzValidateURL -fuzztime=60s ./ims/
func FuzzValidateURL(f *testing.F) {
	f.Add("https://example.com")
	f.Add("http://localhost:8080")
	f.Add("")
	f.Add("not-a-url")
	f.Add("://missing-scheme.com")
	f.Add("https://")

	f.Fuzz(func(t *testing.T, u string) {
		_ = validateURL(u)
	})
}

// TestFuzzDecodeToken generates random inputs for 10 seconds to verify that
// DecodeToken never panics regardless of input. Runs in parallel with other tests.
//
// For deeper exploration, use Go's built-in fuzz engine:
//
//	go test -fuzz=FuzzDecodeToken -fuzztime=60s ./ims/
func TestFuzzDecodeToken(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	deadline := time.After(10 * time.Second)
	iterations := 0

	for {
		select {
		case <-deadline:
			t.Logf("fuzz: %d iterations without panic", iterations)
			return
		default:
			// Generate random JWT-like strings (three dot-separated parts)
			input := randomString(rng, rng.Intn(128)) + "." +
				randomString(rng, rng.Intn(256)) + "." +
				randomString(rng, rng.Intn(128))
			c := Config{Token: input}
			_, _ = c.DecodeToken()
			iterations++
		}
	}
}

// FuzzDecodeToken is a standard Go fuzz target for deeper exploration.
// Run manually: go test -fuzz=FuzzDecodeToken -fuzztime=60s ./ims/
func FuzzDecodeToken(f *testing.F) {
	f.Add("eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.signature")
	f.Add("a.b.c")
	f.Add("...")
	f.Add("")
	f.Add("no-dots-at-all")

	f.Fuzz(func(t *testing.T, token string) {
		c := Config{Token: token}
		_, _ = c.DecodeToken()
	})
}
