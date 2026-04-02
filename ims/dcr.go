// Copyright 2026 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

// Dynamic Client Registration (DCR): POST JSON to IMS /ims/register.

package ims

import (
	"fmt"

	"github.com/adobe/ims-go/ims"
)

const (
	dcrRegisterPath        = "/ims/register"
	maxDCRResponseBodySize = 1 << 20 // 1 MiB
)

func (i Config) validateDCRConfig() error {
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

func (i Config) DCRRegister() (string, error) {
	if err := i.validateDCRConfig(); err != nil {
		return "", fmt.Errorf("invalid parameters for client registration: %w", err)
	}

	c, err := i.newIMSClient()
	if err != nil {
		return "", fmt.Errorf("error creating the IMS client: %w", err)
	}

	resp, err := c.DCR(&ims.DCRRequest{
		ClientName:   i.ClientName,
		RedirectURIs: i.RedirectURIs,
		Scopes:       i.Scopes,
	})
	if err != nil {
		return "", fmt.Errorf("error during client registration: %w", err)
	}

	return string(resp.Body), nil
}
