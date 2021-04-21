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

func exchangeCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exchange",
		Aliases: []string{"exch"},
		Short:   "Exchange an access token for another access token.",
		Long:    "Perform the 'Cluster Access Token Exchange Grant' to request a new token with a new user ID or IMS Org ID.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			resp, err := imsConfig.ClusterExchange()
			if err != nil {
				return fmt.Errorf("error exchanging the access token: %v", err)
			}
			fmt.Println(resp.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "p", "", "IMS client secret.")
	cmd.Flags().StringVarP(&imsConfig.AccessToken, "accessToken", "t", "", "Access token.")
	cmd.Flags().StringVarP(&imsConfig.Organization, "organization", "o", "",
		"IMS Organization for the new token. Can't be used in conjunction with userID.")
	cmd.Flags().StringVarP(&imsConfig.UserID, "userID", "u", "",
		"User ID of the new token. Can't be used in conjunction with organization.")
	cmd.Flags().StringSliceVarP(&imsConfig.Scopes, "scopes", "s", []string{""},
		"Scopes to request in the new token. Subset of the scopes of the original token. Optional value, if no "+
			"scopes are requested the same scopes of the original token will be provided")

	return cmd
}
