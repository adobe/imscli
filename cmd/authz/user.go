// Copyright 2020 Adobe. All rights reserved.
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

func UserCmd(imsConfig *ims.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Negotiate a normal user access token.",
		Long: "Perform the 'Authorization Code Grant Flow' by launching a browser and capturing the returned " +
			"IMS token. Without PKCE only private clients are supported.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			resp, err := imsConfig.AuthorizeUser()
			if err != nil {
				return fmt.Errorf("error in user authorization: %v", err)
			}
			fmt.Println(resp)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "p", "", "IMS client secret.")
	cmd.Flags().StringVarP(&imsConfig.Organization, "organization", "o", "", "IMS Organization.")
	cmd.Flags().StringSliceVarP(&imsConfig.Scopes, "scopes", "s", []string{""}, "Scopes to request.")
	cmd.Flags().IntVarP(&imsConfig.Port, "port", "l", 8888, "Local port to be used by the OAuth Client.")

	return cmd
}
