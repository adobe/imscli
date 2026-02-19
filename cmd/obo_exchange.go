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
		Use:     "obo",
		Aliases: []string{"ob"},
		Short:   "On-Behalf-Of token exchange: get a backend token for a user.",
		Long: `Perform the On-Behalf-Of (OBO) token exchange: exchange a user access token for a new token
suitable for backend-to-backend calls on behalf of that user.

SECURITY: Do NOT send OBO access tokens to frontend clients. OBO tokens are intended only for
backend-to-backend communication. They have a short TTL (e.g. 5 minutes) and the full actor
chain is preserved in the act claim for audit.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			resp, err := imsConfig.OBOExchange()
			if err != nil {
				return fmt.Errorf("error during OBO exchange: %v", err)
			}
			fmt.Println(resp.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "p", "", "IMS client secret.")
	cmd.Flags().StringVarP(&imsConfig.AccessToken, "accessToken", "t", "", "User access token (subject token). Do not use service or impersonation tokens.")
	cmd.Flags().StringSliceVarP(&imsConfig.Scopes, "scopes", "s", []string{""},
		"Scopes to request. Must be within the client's configured scope boundary. Optional.")

	return cmd
}
