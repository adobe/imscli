// Copyright 2021 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

package invalidate

import (
	"fmt"

	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func AccessTokenCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "accessToken",
		Aliases: []string{"acc"},
		Short:   "Invalidate an access token.",
		Long:    "Invalidate an access token.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			err := imsConfig.InvalidateToken()
			if err != nil {
				return fmt.Errorf("error invalidating the access token: %v", err)
			}
			fmt.Println("Token invalidated successfully.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.AccessToken, "accessToken", "t", "", "Access token.")
	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS Client ID.")

	return cmd
}
