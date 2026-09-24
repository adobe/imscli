// Copyright 2026 Adobe. All rights reserved.
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
	"encoding/json"
	"fmt"

	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func promiseCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promise",
		Short: "Use the IMS Promise Token flows.",
		Long: `The promise command performs requests against the IMS /ims/token/v4 endpoint.

This command has no effect by itself, the request needs to be specified as a subcommand.
`,
	}
	cmd.AddCommand(
		accessForPromiseCmd(imsConfig),
		promiseForAccessCmd(imsConfig),
	)
	return cmd
}

func accessForPromiseCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access-for-promise",
		Short: "Exchange an access token for a promise token.",
		Long:  `Perform the IMS Promise Token grant (grant_type=promise) to exchange an authenticating access token for a promise token.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			resp, err := imsConfig.PromiseTokenExchange()
			if err != nil {
				return fmt.Errorf("error during the promise token exchange: %w", err)
			}
			if imsConfig.FullOutput {
				data := struct {
					PromiseToken string `json:"promise_token"`
					TokenType    string `json:"token_type"`
					Scope        string `json:"scope"`
					ExpiresIn    int    `json:"expires_in"`
				}{resp.PromiseToken, resp.TokenType, resp.Scope, resp.ExpiresIn}
				jsonData, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					return fmt.Errorf("error marshalling full JSON response: %w", err)
				}
				fmt.Printf("%s\n", jsonData)
				return nil
			}
			fmt.Println(resp.PromiseToken)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "p", "", "IMS client secret.")
	cmd.Flags().StringVarP(&imsConfig.AccessToken, "accessToken", "t", "", "Authenticating access token (only access tokens are accepted).")
	cmd.Flags().StringVarP(&imsConfig.PromiseDefinitionID, "promiseDefinitionID", "d", "", "Promise definition ID.")
	cmd.Flags().StringSliceVarP(&imsConfig.Scopes, "scopes", "s", nil, "Scopes to request in the promise token.")
	cmd.Flags().BoolVarP(&imsConfig.FullOutput, "fullOutput", "F", false, "Output a JSON with the promise token, token type, scope, and expiration.")

	return cmd
}

func promiseForAccessCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promise-for-access",
		Short: "Exchange a promise token for an access token.",
		Long:  `Perform the IMS Promise Exchange grant (grant_type=promise_exchange) to redeem a promise token for a fresh access token and a rotated promise token.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			resp, err := imsConfig.PromiseExchange()
			if err != nil {
				return fmt.Errorf("error during the promise exchange: %w", err)
			}
			if imsConfig.FullOutput {
				data := struct {
					AccessToken           string `json:"access_token"`
					PromiseToken          string `json:"promise_token"`
					PromiseTokenID        string `json:"promise_token_id"`
					TokenType             string `json:"token_type"`
					Scope                 string `json:"scope"`
					ExpiresIn             int    `json:"expires_in"`
					PromiseTokenExpiresIn int    `json:"promise_token_expires_in"`
				}{
					resp.AccessToken, resp.PromiseToken, resp.PromiseTokenID, resp.TokenType,
					resp.Scope, resp.ExpiresIn, resp.PromiseTokenExpiresIn,
				}
				jsonData, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					return fmt.Errorf("error marshalling full JSON response: %w", err)
				}
				fmt.Printf("%s\n", jsonData)
				return nil
			}
			fmt.Println(resp.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "p", "", "IMS client secret.")
	cmd.Flags().StringVarP(&imsConfig.PromiseToken, "promiseToken", "t", "", "Promise token to exchange.")
	cmd.Flags().StringSliceVarP(&imsConfig.Scopes, "scopes", "s", nil, "Optional scopes to request in the new access token.")
	cmd.Flags().BoolVarP(&imsConfig.FullOutput, "fullOutput", "F", false, "Output a JSON with the access token, rotated promise token, and their expirations.")

	return cmd
}
