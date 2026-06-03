// Copyright 2026 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

package authz

import (
	"fmt"

	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func ImplicitCmd(imsConfig *ims.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "implicit",
		Short: "Negotiate an access token using the OAuth 2.0 implicit grant flow.",
		Long: "Perform the 'Implicit Grant Flow' by launching a browser, completing authentication with IMS, " +
			"and capturing the access token. IMS redirects to a static page (default: " +
			ims.DefaultImplicitRedirectURI + ") that converts the URL fragment to a query string and " +
			"forwards it to the local callback server. Public clients only; no client secret is sent. " +
			"Note that the implicit flow is deprecated in OAuth 2.1 — prefer 'pkce' for new use cases.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			resp, err := imsConfig.AuthorizeImplicit()
			if err != nil {
				return fmt.Errorf("error in implicit authorization: %w", err)
			}
			fmt.Println(resp)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringSliceVarP(&imsConfig.Scopes, "scopes", "s", []string{}, "Scopes to request.")
	cmd.Flags().IntVarP(&imsConfig.Port, "port", "l", 8888, "Local port to be used by the OAuth Client. "+
		"Must match the port that the redirector page sends the browser to (the default redirector pins 8888).")
	cmd.Flags().StringVar(&imsConfig.RedirectURI, "redirectURI", ims.DefaultImplicitRedirectURI,
		"Redirect URI registered with IMS. Defaults to the canonical public redirector; override to use a self-hosted page (e.g., when running on a non-default port).")
	cmd.Flags().StringSliceVarP(&imsConfig.Resource, "resource", "r", nil,
		"RFC 8707 resource indicator URI(s) for audience-restricted tokens.")

	return cmd
}
