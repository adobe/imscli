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

	"github.com/adobe/imscli/cmd/pretty"
	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func organizationsCmd(imsConfig *ims.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "organizations",
		Aliases: []string{"orgs"},
		Short:   "Requests the user organizations.",
		Long:    "Requests the user organizations associated to the provided access token.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			resp, err := imsConfig.GetOrganizations()
			if err != nil {
				return fmt.Errorf("error in get organizations cmd: %w", err)
			}
			fmt.Println(pretty.JSON(resp))
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.AccessToken, "accessToken", "t", "", "Access token.")
	cmd.Flags().StringVarP(&imsConfig.OrgsApiVersion, "orgsApiVersion", "a", "v5", "Organizations API version.")

	return cmd
}
