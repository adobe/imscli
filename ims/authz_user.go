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
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/adobe/ims-go/ims"
	"github.com/adobe/ims-go/login"
	"github.com/pkg/browser"
)

// Validate that:
//   - the ims.Config struct has the necessary parameters for AuthorizeUser
//   - the provided environment exists
func (i Config) validateAuthorizeUserConfig() error {

	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case !validateURL(i.URL):
		return fmt.Errorf("unable to parse URL parameter")
	case len(i.Scopes) == 0 || i.Scopes[0] == "":
		return fmt.Errorf("missing scopes parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client id parameter")
	case i.Organization == "":
		return fmt.Errorf("missing organization parameter")
	case i.Port <= 0:
		return fmt.Errorf("missing or invalid port parameter")
	case i.ClientSecret == "":
		if i.PublicClient {
			log.Println("all needed parameters verified not empty")
			return nil
		}
		return fmt.Errorf("missing client secret parameter")
	default:
		log.Println("all needed parameters verified not empty")
	}

	return nil
}

// AuthorizeUser uses the standard OAuth2 authorization code grant flow.
func (i Config) AuthorizeUser() (string, error) {
	return i.authorizeUser(false)
}

// AuthorizeUserPKCE uses the OAuth2 authorization code grant flow with PKCE.
func (i Config) AuthorizeUserPKCE() (string, error) {
	return i.authorizeUser(true)
}

func (i Config) authorizeUser(pkce bool) (string, error) {
	// Perform parameter validation
	err := i.validateAuthorizeUserConfig()
	if err != nil {
		return "", fmt.Errorf("invalid parameters for login user: %w", err)
	}

	c, err := i.newIMSClient()
	if err != nil {
		return "", fmt.Errorf("error creating the IMS client: %w", err)
	}

	server, err := login.NewServer(&login.ServerConfig{
		Client:       c,
		ClientID:     i.ClientID,
		ClientSecret: i.ClientSecret,
		Scope:        i.Scopes,
		UsePKCE:      pkce,
		RedirectURI:  fmt.Sprintf("http://localhost:%d", i.Port),
		OnError: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `
				<h1>An error occurred</h1>
				<p>Please look at the terminal output for further details.</p>
			`)
		}),
		OnSuccess: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `
				<h1>Login successful!</h1>
				<p>You can close this tab.</p>
			`)
		}),
	})
	if err != nil {
		return "", fmt.Errorf("create authorization server: %w", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", i.Port))
	if err != nil {
		return "", fmt.Errorf("unable to listen at port %d", i.Port)
	}
	defer listener.Close()

	log.Println("Local server successfully launched and contacted.")

	localUrl := fmt.Sprintf("http://localhost:%d/", i.Port)

	// Suppress "Opening in existing browser session." messages from chromium-based
	// browsers. The CLI token output goes to stdout, so stray browser messages
	// would corrupt piped/scripted output. Save and restore to avoid permanent
	// mutation of the package-level variable.
	origStdout := browser.Stdout
	browser.Stdout = nil
	err = browser.OpenURL(localUrl)
	browser.Stdout = origStdout
	if err != nil {
		fmt.Fprintf(os.Stderr, "error launching the browser, open it and visit %s\n", localUrl)
	}

	// Capture Serve errors via a buffered channel. Buffered so the goroutine
	// can always write and exit, even if nobody reads (e.g., a response arrived
	// first). See docs/oauth-serve-error.md for a detailed explanation.
	serveCh := make(chan error, 1)
	go func() {
		serveCh <- server.Serve(listener)
	}()

	var (
		serr error
		resp *ims.TokenResponse
	)

	select {
	case serr = <-server.Error():
		log.Println("The IMS HTTP handler returned an error message.")
	case resp = <-server.Response():
		log.Println("The IMS HTTP handler returned a message.")
	case serr = <-serveCh:
		log.Println("The local server stopped unexpectedly.")
	case <-time.After(time.Minute * 5):
		fmt.Fprintf(os.Stderr, "Timeout reached waiting for the user to finish the authentication ...\n")
		serr = fmt.Errorf("user timed out")
	}

	// Drain channels to prevent a deadlock between Shutdown() waiting for
	// in-flight handlers and handlers blocking on unbuffered channel writes.
	// See docs/oauth-shutdown-deadlock.md for a detailed explanation.
	go func() {
		for range server.Response() {
		}
	}()
	go func() {
		for range server.Error() {
		}
	}()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		return "", fmt.Errorf("error shutting down the local server: %w", err)
	}
	log.Println("Local server shut down ...")

	if serr != nil {
		return "", fmt.Errorf("error negotiating the authorization code: %w", serr)
	}
	log.Println("No error from Authorization Code handler, server is successfully shut down.")

	return resp.AccessToken, nil
}
