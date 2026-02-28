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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (i Config) validateRegisterConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case !validateURL(i.URL):
		return fmt.Errorf("invalid IMS base URL parameter")
	case i.ClientName == "":
		return fmt.Errorf("missing client name parameter")
	case len(i.RedirectURIs) == 0:
		return fmt.Errorf("missing redirect URIs parameter")
	default:
		return nil
	}
}

// Register performs Dynamic Client Registration.
func (i Config) Register() (string, error) {
	if err := i.validateRegisterConfig(); err != nil {
		return "", fmt.Errorf("invalid parameters for client registration: %w", err)
	}

	// Build the request payload using json.Marshal for proper escaping.
	payload, err := json.Marshal(map[string]interface{}{
		"client_name":   i.ClientName,
		"redirect_uris": i.RedirectURIs,
	})
	if err != nil {
		return "", fmt.Errorf("error building registration payload: %w", err)
	}

	endpoint := strings.TrimRight(i.URL, "/") + "/ims/register"

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(string(payload)))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("content-type", "application/json")

	httpClient, err := i.httpClient()
	if err != nil {
		return "", fmt.Errorf("error creating the HTTP client: %w", err)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making registration request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}
