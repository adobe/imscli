// Copyright 2020 Adobe. All rights reserved.
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

func authzJWTCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jwt",
		Short: "Negotiate a JWT token",
		Long:  `Perform the 'Assertion Grant Type Flow' to request a token.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			resp, err := imsConfig.AuthorizeJWTExchange()
			if err != nil {
				return fmt.Errorf("error in jwt authorization: %v", err)
			}
			fmt.Println(resp.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS Client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "s", "", "IMS Client secret.")
	cmd.Flags().StringVarP(&imsConfig.Organization, "organization", "o", "", "IMS Organization.")
	cmd.Flags().StringVarP(&imsConfig.Account, "account", "a", "", "Technical Account ID.")
	cmd.Flags().StringVarP(&imsConfig.PrivateKeyPath, "privateKey", "k", "", "Private Key file.")
	cmd.Flags().StringSliceVarP(&imsConfig.Metascopes, "metascopes", "m", []string{""}, "Metascopes to request.")

	return cmd
}
