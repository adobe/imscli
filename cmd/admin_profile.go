// Copyright 2023 Adobe. All rights reserved.
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

func adminProfileCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Requests the user profile using the admin API.",
		Long:  "Requests the specified user profile using the admin API and a service token.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			resp, err := imsConfig.GetProfile()
			if err != nil {
				return fmt.Errorf("error in get profile cmd: %v", err)
			}
			fmt.Println(resp)
			return nil
		},
	}
	cmd.Flags().StringVarP(&imsConfig.ServiceToken, "serviceToken", "t", "", "Service token.")
	cmd.Flags().StringVarP(&imsConfig.ProfileApiVersion, "profileApiVersion", "a", "v1", "Profile API version.")

	return cmd
}
