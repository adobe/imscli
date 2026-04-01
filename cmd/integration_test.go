// Copyright 2024 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ---------- helpers ----------

type capturedRequest struct {
	Method string
	Path   string
	Header http.Header
	Form   map[string]string
	Query  map[string]string
}

type requestLog struct {
	received bool
	capturedRequest
}

func (r *requestLog) record(req *http.Request) {
	_ = req.ParseForm()

	form := make(map[string]string)
	for k, v := range req.PostForm {
		form[k] = v[0]
	}
	query := make(map[string]string)
	for k, v := range req.URL.Query() {
		query[k] = v[0]
	}

	r.received = true
	r.Method = req.Method
	r.Path = req.URL.Path
	r.Header = req.Header.Clone()
	r.Form = form
	r.Query = query
}

func newMockIMS(t *testing.T) (*httptest.Server, *requestLog) {
	t.Helper()
	rlog := &requestLog{}

	mux := http.NewServeMux()

	handle := func(pattern string, body string) {
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			rlog.record(r)
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, body)
		})
	}

	handle("POST /ims/validate_token/v1", `{"valid":true}`)
	handle("POST /ims/invalidate_token/v2", `{}`)

	// Profile and organizations endpoints use versioned paths.
	// Use prefix-matching by trailing slash for version flexibility.
	mux.HandleFunc("/ims/profile/", func(w http.ResponseWriter, r *http.Request) {
		rlog.record(r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"name":"test-user"}`)
	})
	mux.HandleFunc("/ims/organizations/", func(w http.ResponseWriter, r *http.Request) {
		rlog.record(r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `[{"orgName":"test-org"}]`)
	})

	handle("POST /ims/token/v2", `{"access_token":"at","refresh_token":"rt","expires_in":3600}`)
	handle("POST /ims/token/v3", `{"access_token":"at","expires_in":3600}`)

	mux.HandleFunc("/ims/admin_profile/", func(w http.ResponseWriter, r *http.Request) {
		rlog.record(r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"name":"admin"}`)
	})
	mux.HandleFunc("/ims/admin_organizations/", func(w http.ResponseWriter, r *http.Request) {
		rlog.record(r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `[{"orgName":"admin-org"}]`)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, rlog
}

func execCmd(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	cmd := RootCmd("test")
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(outBuf)
	cmd.SetErr(errBuf)
	cmd.SetArgs(args)
	err = cmd.Execute()
	return outBuf.String(), errBuf.String(), err
}

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "imscli.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// ---------- 1. Precedence ----------

func TestPrecedence_FlagOverridesEnv(t *testing.T) {
	srv, rlog := newMockIMS(t)
	t.Setenv("IMS_URL", "http://invalid.example.com")
	_, _, err := execCmd(t, "validate", "accessToken",
		"--url", srv.URL,
		"--clientID", "cid",
		"--accessToken", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rlog.received {
		t.Fatal("mock server received no requests — flag did not override env")
	}
}

func TestPrecedence_EnvOverridesConfigFile(t *testing.T) {
	srv, rlog := newMockIMS(t)
	cfg := writeConfigFile(t, "url: http://invalid.example.com\nclientid: cid\naccesstoken: tok\n")
	t.Setenv("IMS_URL", srv.URL)
	_, _, err := execCmd(t, "validate", "accessToken",
		"--configFile", cfg,
		"--clientID", "cid",
		"--accessToken", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rlog.received {
		t.Fatal("mock server received no requests — env did not override config file")
	}
}

func TestPrecedence_ConfigFileOverridesDefault(t *testing.T) {
	srv, rlog := newMockIMS(t)
	cfg := writeConfigFile(t, "url: "+srv.URL+"\nclientid: cid\naccesstoken: tok\n")
	_, _, err := execCmd(t, "validate", "accessToken",
		"--configFile", cfg,
		"--clientID", "cid",
		"--accessToken", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rlog.received {
		t.Fatal("mock server received no requests — config file URL not used")
	}
}

func TestPrecedence_DefaultUsedWhenNothingSet(t *testing.T) {
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "validate", "accessToken",
		"--configFile", empty,
		"--clientID", "cid",
		"--accessToken", "tok")
	if err == nil {
		t.Fatal("expected an error when using default URL (unreachable)")
	}
}

// ---------- 2. URL reaches server ----------

func TestURL_SubcommandRouting(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		path   string
		method string
	}{
		{
			name:   "validate accessToken",
			args:   []string{"validate", "accessToken", "--clientID", "cid", "--accessToken", "tok"},
			path:   "/ims/validate_token/v1",
			method: "POST",
		},
		{
			name:   "invalidate accessToken",
			args:   []string{"invalidate", "accessToken", "--clientID", "cid", "--accessToken", "tok"},
			path:   "/ims/invalidate_token/v2",
			method: "POST",
		},
		{
			name:   "profile",
			args:   []string{"profile", "--accessToken", "tok"},
			path:   "/ims/profile/v1",
			method: "GET",
		},
		{
			name:   "organizations",
			args:   []string{"organizations", "--accessToken", "tok"},
			path:   "/ims/organizations/v5",
			method: "GET",
		},
		{
			name:   "refresh",
			args:   []string{"refresh", "--clientID", "cid", "--clientSecret", "sec", "--refreshToken", "rt"},
			path:   "/ims/token/v2",
			method: "POST",
		},
		{
			name:   "exchange",
			args:   []string{"exchange", "--clientID", "cid", "--clientSecret", "sec", "--accessToken", "tok", "--organization", "org"},
			path:   "/ims/token/v3",
			method: "POST",
		},
		{
			name:   "admin profile",
			args:   []string{"admin", "profile", "--clientID", "cid", "--serviceToken", "st", "--guid", "g1", "--authSrc", "as1"},
			path:   "/ims/admin_profile/v1",
			method: "POST",
		},
		{
			name:   "admin organizations",
			args:   []string{"admin", "organizations", "--clientID", "cid", "--serviceToken", "st", "--guid", "g1", "--authSrc", "as1"},
			path:   "/ims/admin_organizations/v5",
			method: "POST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, rlog := newMockIMS(t)
			empty := writeConfigFile(t, "")
			args := append([]string{"--url", srv.URL, "--configFile", empty}, tt.args...)
			_, _, err := execCmd(t, args...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !rlog.received {
				t.Fatal("mock server received no requests")
			}
			got := rlog.capturedRequest
			if got.Path != tt.path {
				t.Errorf("path = %q, want %q", got.Path, tt.path)
			}
			if got.Method != tt.method {
				t.Errorf("method = %q, want %q", got.Method, tt.method)
			}
		})
	}
}

// ---------- 3. Timeout ----------

func TestTimeout_FlagTimeout(t *testing.T) {
	done := make(chan struct{})
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(10 * time.Second):
		case <-done:
		}
	}))
	t.Cleanup(func() {
		close(done)
		slow.Close()
	})
	empty := writeConfigFile(t, "")

	start := time.Now()
	_, _, err := execCmd(t, "validate", "accessToken",
		"--url", slow.URL,
		"--configFile", empty,
		"--timeout", "1",
		"--clientID", "cid",
		"--accessToken", "tok")
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed > 3*time.Second {
		t.Errorf("took %v, expected ~1s timeout", elapsed)
	}
}

func TestTimeout_EnvTimeout(t *testing.T) {
	done := make(chan struct{})
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(10 * time.Second):
		case <-done:
		}
	}))
	t.Cleanup(func() {
		close(done)
		slow.Close()
	})
	empty := writeConfigFile(t, "")

	t.Setenv("IMS_TIMEOUT", "1")
	start := time.Now()
	_, _, err := execCmd(t, "validate", "accessToken",
		"--url", slow.URL,
		"--configFile", empty,
		"--clientID", "cid",
		"--accessToken", "tok")
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed > 3*time.Second {
		t.Errorf("took %v, expected ~1s timeout", elapsed)
	}
}

// ---------- 4. Config file loading ----------

func TestConfigFile_AllValuesFromFile(t *testing.T) {
	srv, rlog := newMockIMS(t)
	cfg := writeConfigFile(t, "url: "+srv.URL+"\nclientid: file-cid\naccesstoken: file-tok\n")
	_, _, err := execCmd(t, "validate", "accessToken",
		"--configFile", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rlog.received {
		t.Fatal("mock server received no requests")
	}
	got := rlog.capturedRequest
	if got.Form["client_id"] != "file-cid" {
		t.Errorf("client_id = %q, want %q", got.Form["client_id"], "file-cid")
	}
}

func TestConfigFile_ExplicitPath(t *testing.T) {
	srv, rlog := newMockIMS(t)
	dir := t.TempDir()
	p := filepath.Join(dir, "custom-name.yaml")
	content := "url: " + srv.URL + "\nclientid: cid\naccesstoken: tok\n"
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	_, _, err := execCmd(t, "validate", "accessToken", "--configFile", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rlog.received {
		t.Fatal("mock server received no requests")
	}
}

func TestConfigFile_InvalidPath(t *testing.T) {
	_, _, err := execCmd(t, "validate", "accessToken",
		"--configFile", "/nonexistent/path/imscli.yaml",
		"--clientID", "cid",
		"--accessToken", "tok")
	if err == nil {
		t.Fatal("expected error for invalid config file path")
	}
	if !strings.Contains(err.Error(), "unable to read configuration file") {
		t.Errorf("error = %q, want to contain 'unable to read configuration file'", err.Error())
	}
}

// ---------- 5. Command-specific flags ----------

func TestCommandSpecific_ClientIDInForm(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "validate", "accessToken",
		"--url", srv.URL, "--configFile", empty,
		"--clientID", "my-client",
		"--accessToken", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	if got.Form["client_id"] != "my-client" {
		t.Errorf("client_id = %q, want %q", got.Form["client_id"], "my-client")
	}
}

func TestCommandSpecific_XImsClientIdHeader(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "validate", "accessToken",
		"--url", srv.URL, "--configFile", empty,
		"--clientID", "my-client",
		"--accessToken", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	if got.Header.Get("X-Ims-Clientid") != "my-client" {
		t.Errorf("X-IMS-ClientId = %q, want %q", got.Header.Get("X-Ims-Clientid"), "my-client")
	}
}

func TestCommandSpecific_AuthorizationHeader(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "profile",
		"--url", srv.URL, "--configFile", empty,
		"--accessToken", "my-access-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	want := "Bearer my-access-token"
	if got.Header.Get("Authorization") != want {
		t.Errorf("Authorization = %q, want %q", got.Header.Get("Authorization"), want)
	}
}

func TestCommandSpecific_CascadingFlag(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "invalidate", "refreshToken",
		"--url", srv.URL, "--configFile", empty,
		"--clientID", "cid",
		"--refreshToken", "rt",
		"--cascading")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	if got.Form["cascading"] != "all" {
		t.Errorf("cascading = %q, want %q", got.Form["cascading"], "all")
	}
}

func TestCommandSpecific_ClientSecretInForm(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "invalidate", "serviceToken",
		"--url", srv.URL, "--configFile", empty,
		"--clientID", "cid",
		"--serviceToken", "st",
		"--clientSecret", "my-secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	if got.Form["client_secret"] != "my-secret" {
		t.Errorf("client_secret = %q, want %q", got.Form["client_secret"], "my-secret")
	}
}

func TestCommandSpecific_ExchangeClientIDInQuery(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "exchange",
		"--url", srv.URL, "--configFile", empty,
		"--clientID", "exch-cid",
		"--clientSecret", "sec",
		"--accessToken", "tok",
		"--organization", "org1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	if got.Query["client_id"] != "exch-cid" {
		t.Errorf("query client_id = %q, want %q", got.Query["client_id"], "exch-cid")
	}
}

func TestCommandSpecific_RefreshFormData(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "refresh",
		"--url", srv.URL, "--configFile", empty,
		"--clientID", "cid",
		"--clientSecret", "sec",
		"--refreshToken", "rt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	if got.Form["grant_type"] != "refresh_token" {
		t.Errorf("grant_type = %q, want %q", got.Form["grant_type"], "refresh_token")
	}
}

func TestCommandSpecific_AdminGuidAuthSrc(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "admin", "profile",
		"--url", srv.URL, "--configFile", empty,
		"--clientID", "cid",
		"--serviceToken", "st",
		"--guid", "user-guid",
		"--authSrc", "my-auth-src")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	if got.Form["guid"] != "user-guid" {
		t.Errorf("guid = %q, want %q", got.Form["guid"], "user-guid")
	}
	if got.Form["auth_src"] != "my-auth-src" {
		t.Errorf("auth_src = %q, want %q", got.Form["auth_src"], "my-auth-src")
	}
}

func TestCommandSpecific_ClientCredentialsOrgID(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "authorize", "clientCredentials",
		"--url", srv.URL, "--configFile", empty,
		"--clientID", "cid",
		"--clientSecret", "sec",
		"--scopes", "openid",
		"--organization", "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	if got.Form["org_id"] != "my-org" {
		t.Errorf("org_id = %q, want %q", got.Form["org_id"], "my-org")
	}
}

func TestCommandSpecific_ClientCredentialsNoOrgID(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	_, _, err := execCmd(t, "authorize", "clientCredentials",
		"--url", srv.URL, "--configFile", empty,
		"--clientID", "cid",
		"--clientSecret", "sec",
		"--scopes", "openid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	if _, ok := got.Form["org_id"]; ok {
		t.Errorf("org_id should not be present in request, got %q", got.Form["org_id"])
	}
}

// ---------- 6. API version flags ----------

func TestAPIVersion_Routing(t *testing.T) {
	tests := []struct {
		name string
		args []string
		path string
	}{
		{
			name: "profile default v1",
			args: []string{"profile", "--accessToken", "tok"},
			path: "/ims/profile/v1",
		},
		{
			name: "profile v2",
			args: []string{"profile", "--accessToken", "tok", "--profileApiVersion", "v2"},
			path: "/ims/profile/v2",
		},
		{
			name: "profile v3",
			args: []string{"profile", "--accessToken", "tok", "--profileApiVersion", "v3"},
			path: "/ims/profile/v3",
		},
		{
			name: "organizations default v5",
			args: []string{"organizations", "--accessToken", "tok"},
			path: "/ims/organizations/v5",
		},
		{
			name: "organizations v6",
			args: []string{"organizations", "--accessToken", "tok", "--orgsApiVersion", "v6"},
			path: "/ims/organizations/v6",
		},
		{
			name: "admin profile v2",
			args: []string{"admin", "profile", "--clientID", "cid", "--serviceToken", "st", "--guid", "g", "--authSrc", "a", "--profileApiVersion", "v2"},
			path: "/ims/admin_profile/v2",
		},
		{
			name: "admin organizations v6",
			args: []string{"admin", "organizations", "--clientID", "cid", "--serviceToken", "st", "--guid", "g", "--authSrc", "a", "--orgsApiVersion", "v6"},
			path: "/ims/admin_organizations/v6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, rlog := newMockIMS(t)
			empty := writeConfigFile(t, "")
			args := append([]string{"--url", srv.URL, "--configFile", empty}, tt.args...)
			_, _, err := execCmd(t, args...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !rlog.received {
				t.Fatal("mock server received no requests")
			}
			got := rlog.capturedRequest
			if got.Path != tt.path {
				t.Errorf("path = %q, want %q", got.Path, tt.path)
			}
		})
	}
}

// ---------- 7. Env var mapping ----------

func TestEnvVar_ClientID(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	t.Setenv("IMS_CLIENTID", "env-cid")
	_, _, err := execCmd(t, "validate", "accessToken",
		"--url", srv.URL, "--configFile", empty,
		"--accessToken", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	if got.Form["client_id"] != "env-cid" {
		t.Errorf("client_id = %q, want %q", got.Form["client_id"], "env-cid")
	}
}

func TestEnvVar_AccessToken(t *testing.T) {
	srv, rlog := newMockIMS(t)
	empty := writeConfigFile(t, "")
	t.Setenv("IMS_ACCESSTOKEN", "env-token")
	_, _, err := execCmd(t, "profile",
		"--url", srv.URL, "--configFile", empty)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := rlog.capturedRequest
	want := "Bearer env-token"
	if got.Header.Get("Authorization") != want {
		t.Errorf("Authorization = %q, want %q", got.Header.Get("Authorization"), want)
	}
}
