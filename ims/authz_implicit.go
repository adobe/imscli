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
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/adobe/ims-go/ims"
	"github.com/pkg/browser"
)

// DefaultImplicitRedirectURI is the canonical public redirector served from
// docs/redirect/implicit/index.html. The implicit flow returns the access
// token in the URL fragment (#access_token=...), which the browser never
// sends to a server. That static page runs a one-liner that converts the
// fragment to a query string and navigates the browser to localhost, where
// our Go handler can finally read the parameters from r.URL.Query().
//
// Crafty users running on a non-default port can host their own redirector
// (pointing at the desired localhost port), register it with IMS for their
// client, and pass --redirectURI to override this default.
const DefaultImplicitRedirectURI = "https://opensource.adobe.com/imscli/redirect/implicit/"

// validateAuthorizeImplicitConfig checks that the configuration has valid
// parameters for the implicit grant flow.
func (i Config) validateAuthorizeImplicitConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case !validateURL(i.URL):
		return fmt.Errorf("unable to parse URL parameter")
	case len(i.Scopes) == 0 || i.Scopes[0] == "":
		return fmt.Errorf("missing scopes parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client id parameter")
	case i.Port <= 0:
		return fmt.Errorf("missing or invalid port parameter")
	case i.RedirectURI == "":
		return fmt.Errorf("missing redirect URI parameter")
	case !validateURL(i.RedirectURI):
		return fmt.Errorf("unable to parse redirect URI parameter")
	}
	log.Println("all needed parameters verified not empty")
	return nil
}

// AuthorizeImplicit performs the OAuth 2.0 implicit grant flow with IMS.
// IMS redirects the browser to the configured RedirectURI (a public static
// page) which JS-rewrites the URL fragment into a query string and navigates
// the browser to the local listener. Returns the access token after state
// validation.
func (i Config) AuthorizeImplicit() (string, error) {
	if err := i.validateAuthorizeImplicitConfig(); err != nil {
		return "", fmt.Errorf("invalid parameters for implicit authorization: %w", err)
	}

	c, err := i.newIMSClient()
	if err != nil {
		return "", fmt.Errorf("error creating the IMS client: %w", err)
	}

	state, err := randomState()
	if err != nil {
		return "", fmt.Errorf("generate state: %w", err)
	}

	authURL, err := c.AuthorizeURL(&ims.AuthorizeURLConfig{
		ClientID:    i.ClientID,
		GrantType:   ims.GrantTypeImplicit,
		Scope:       i.Scopes,
		RedirectURI: i.RedirectURI,
		State:       state,
		Resource:    i.Resource,
	})
	if err != nil {
		return "", fmt.Errorf("build authorize URL: %w", err)
	}

	// Buffered channels: handlers can always send and exit, even if the main
	// goroutine has already moved on (e.g., after a timeout). Avoids the
	// shutdown deadlock pattern documented in docs/oauth-shutdown-deadlock.md.
	resCh := make(chan *TokenInfo, 1)
	errCh := make(chan error, 1)

	handler := &implicitHandler{
		expectedState: state,
		resCh:         resCh,
		errCh:         errCh,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.capture)

	server := &http.Server{Handler: mux}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", i.Port))
	if err != nil {
		return "", fmt.Errorf("unable to listen at port %d", i.Port)
	}
	defer func() { _ = listener.Close() }()

	log.Println("Local server successfully launched and contacted.")

	// Suppress chromium "Opening in existing browser session" messages; the CLI
	// token output goes to stdout, so stray browser messages would corrupt
	// piped/scripted output. Save and restore to avoid permanent mutation of
	// the package-level variable.
	origStdout := browser.Stdout
	browser.Stdout = nil
	err = browser.OpenURL(authURL)
	browser.Stdout = origStdout
	if err != nil {
		fmt.Fprintf(os.Stderr, "error launching the browser, open it and visit %s\n", authURL)
	}

	// Capture Serve errors via a buffered channel. See ims/authz_user.go:136-140
	// for the rationale.
	serveCh := make(chan error, 1)
	go func() {
		serveCh <- server.Serve(listener)
	}()

	var (
		serr error
		resp *TokenInfo
	)

	select {
	case serr = <-errCh:
		log.Println("The implicit callback handler returned an error message.")
	case resp = <-resCh:
		log.Println("The implicit callback handler returned a token.")
	case serr = <-serveCh:
		log.Println("The local server stopped unexpectedly.")
	case <-time.After(authTimeout):
		fmt.Fprintf(os.Stderr, "Timeout reached waiting for the user to finish the authentication ...\n")
		serr = fmt.Errorf("user timed out")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return "", fmt.Errorf("error shutting down the local server: %w", err)
	}
	log.Println("Local server shut down ...")

	if serr != nil {
		return "", fmt.Errorf("error in implicit authorization: %w", serr)
	}

	return resp.AccessToken, nil
}

// implicitHandler holds the per-request state for the implicit-flow callback.
// Extracted to a struct so the route handler can be unit-tested in isolation
// from net.Listen and the live select loop.
type implicitHandler struct {
	expectedState string
	resCh         chan<- *TokenInfo
	errCh         chan<- error
}

// capture reads the access token (and related params) from the query string
// that the external redirector page rewrote from the URL fragment. Validates
// state against the value sent on the original authorize request — protects
// the local callback from CSRF (any other open browser tab could otherwise
// navigate here with attacker-supplied parameters).
func (h *implicitHandler) capture(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	got := q.Get("state")
	if subtle.ConstantTimeCompare([]byte(got), []byte(h.expectedState)) != 1 {
		h.errCh <- fmt.Errorf("state mismatch")
		writeCallbackHTML(w, `<h1>Login failed</h1><p>State mismatch — see terminal.</p>`)
		return
	}

	if e := q.Get("error"); e != "" {
		h.errCh <- fmt.Errorf("authorization error: %s: %s", e, q.Get("error_description"))
		writeCallbackHTML(w, `<h1>Login failed</h1><p>See terminal.</p>`)
		return
	}

	token := q.Get("access_token")
	if token == "" {
		h.errCh <- fmt.Errorf("missing access_token in callback")
		writeCallbackHTML(w, `<h1>Login failed</h1><p>No token in response.</p>`)
		return
	}

	h.resCh <- &TokenInfo{AccessToken: token}
	writeCallbackHTML(w, `<h1>Login successful!</h1><p>You can close this tab.</p>`)
}

func writeCallbackHTML(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.WriteString(w, body)
}

// randomState generates a cryptographically random state parameter for the
// authorize request. Mirrors github.com/adobe/ims-go/login/server.go.
func randomState() (string, error) {
	b := make([]byte, 128)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random state: %w", err)
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
