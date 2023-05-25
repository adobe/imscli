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

const port = 8888

/*
 * Validate that:
 *	- the ims.Config struct has the necessary parameters for AuthorizeUser
 *  - the provided environment exists
 */
func (i Config) validateAuthorizeUserConfig() error {

	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case !validateURL(i.URL):
		return fmt.Errorf("unable to parse URL parameter")
	case i.Scopes[0] == "":
		return fmt.Errorf("missing scopes parameter")
	case i.ClientID == "":
		return fmt.Errorf("missing client id parameter")
	case i.ClientSecret == "":
		return fmt.Errorf("missing client secret parameter")
	case i.Organization == "":
		return fmt.Errorf("missing organization parameter")
	default:
		log.Println("all needed parameters verified not empty")
	}

	return nil
}

/*
 * AuthorizeUser uses the standard Oauth2 authorization code grant flow. The Oauth2 configuration is
 * taken from the Config struct.
 */
func (i Config) AuthorizeUser() (string, error) {
	// Perform parameter validation
	err := i.validateAuthorizeUserConfig()
	if err != nil {
		return "", fmt.Errorf("invalid parameters for login user: %v", err)
	}

	httpClient, err := i.httpClient()
	if err != nil {
		return "", fmt.Errorf("error creating the HTTP Client: %v", err)
	}

	c, err := ims.NewClient(&ims.ClientConfig{
		URL:    i.URL,
		Client: httpClient,
	})
	if err != nil {
		return "", fmt.Errorf("error during client creation: %v", err)
	}

	server, err := login.NewServer(&login.ServerConfig{
		Client:       c,
		ClientID:     i.ClientID,
		ClientSecret: i.ClientSecret,
		Scope:        i.Scopes,
		RedirectURI:  "http://localhost:8888",
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
		return "", fmt.Errorf("create authorization server: %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return "", fmt.Errorf("unable to listen at port %d", port)
	}

	log.Println("Local server successfully launched and contacted.")

	localUrl := fmt.Sprintf("http://localhost:%d/", port)

	// redirect stdout to avoid "Opening in existing browser session." message from chromium
	browser.Stdout = nil

	err = browser.OpenURL(localUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error launching the browser, open it and visit %s\n", localUrl)
	}

	go server.Serve(listener)

	var (
		serr error
		resp *ims.TokenResponse
	)

	select {
	case serr = <-server.Error():
		log.Println("The IMS HTTP handler returned a message.")
	case resp = <-server.Response():
		log.Println("The IMS HTTP handler returned a message.")
	case <-time.After(time.Minute * 5):
		fmt.Fprintf(os.Stderr, "Timeout reached waiting for the user to finish the authentication ...\n")
		serr = fmt.Errorf("user timed out")
	}

	if err = server.Shutdown(context.Background()); err != nil {
		return "", fmt.Errorf("error shutting down the local server: %v", err)
	}
	log.Println("Local server shut down ...")

	if serr != nil {
		return "", fmt.Errorf("error negotiating the authorization code: %v", serr)
	}
	log.Println("No error from Authorization Code handler, server is successfully shut down.")

	return resp.AccessToken, nil
}
