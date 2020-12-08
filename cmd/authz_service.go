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

func authzServiceCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Negotiate authorization as a service",
		Long:  `Perform the 'Client Credential Authorization Flow' to negotiate an access token for a service.'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			resp, err := imsConfig.AuthorizeService()
			if err != nil {
				return fmt.Errorf("error in login service: %v", err)
			}
			fmt.Println(resp)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "p", "", "IMS client secret.")
	cmd.Flags().StringVarP(&imsConfig.ServiceToken, "serviceToken", "s", "", "Service token.")

	return cmd
}
