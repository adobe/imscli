// Copyright 2021 Adobe. All rights reserved.
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
	"fmt"

	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func oboExchangeCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "on-behalf-of",
		Aliases: []string{"obo"},
		Short:   "On-Behalf-Of token exchange.",
		Long:    `Token exchange using the RFC 8693 (urn:ietf:params:oauth:grant-type:token-exchange) grant.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			resp, err := imsConfig.OBOExchange()
			if err != nil {
				return fmt.Errorf("error during On-Behalf-Of exchange: %w", err)
			}
			fmt.Println(resp.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "p", "", "IMS client secret.")
	cmd.Flags().StringVarP(&imsConfig.AccessToken, "accessToken", "t", "", "User access token (subject token). Only access tokens are accepted.")
	cmd.Flags().StringSliceVarP(&imsConfig.Scopes, "scopes", "s", nil,
		"Optional scopes to request; if omitted, none are sent. When set, must stay within the client's configured scope boundary.")
	cmd.Flags().StringSliceVarP(&imsConfig.Resource, "resource", "r", nil, "RFC 8707 resource indicator URI(s) for audience-restricted tokens.")

	return cmd
}
