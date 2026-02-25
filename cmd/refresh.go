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
	"encoding/json"
	"fmt"

	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func refreshCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "refresh",
		Aliases: []string{"ref"},
		Short:   "Exchange a refresh token for new access and refresh tokens.",
		Long:    "Exchange a refresh token for new access and refresh tokens.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			resp, err := imsConfig.Refresh()
			if err != nil {
				return fmt.Errorf("error during the token refresh: %v", err)
			}
			if imsConfig.FullOutput {
				data := map[string]interface{}{
					"access_token":  resp.AccessToken,
					"refresh_token": resp.RefreshToken,
				}
				jsonData, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					return fmt.Errorf("error marshalling full JSON response: %v", err)
				}
				fmt.Printf("%s", jsonData)
				return nil
			}
			fmt.Println(resp.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "p", "", "IMS client secret.")
	cmd.Flags().StringVarP(&imsConfig.RefreshToken, "refreshToken", "t", "", "Refresh token.")
	cmd.Flags().StringSliceVarP(&imsConfig.Scopes, "scopes", "s", []string{},
		"Scopes to request in the new token. Subset of the scopes of the original token. Optional value, if no "+
			"scopes are requested the same scopes of the original token will be provided.")
	cmd.Flags().BoolVarP(&imsConfig.FullOutput, "fullOutput", "F", false, "Output a JSON with access and refresh tokens.")

	return cmd
}
