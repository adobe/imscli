// Copyright 2025 Adobe. All rights reserved.
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

	"github.com/adobe/imscli/cmd/prettify"
	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func dcrCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dcr",
		Short: "Dynamic Client Registration operations.",
		Long:  `The dcr command enables Dynamic Client Registration operations.`,
	}
	cmd.AddCommand(
		registerCmd(imsConfig),
	)
	return cmd
}

func registerCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a client.",
		Long:  `Register a new OAuth client using Dynamic Client Registration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			resp, err := imsConfig.Register()
			if err != nil {
				return fmt.Errorf("error during client registration: %w", err)
			}

			fmt.Println(prettify.JSON(resp))
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientName, "clientName", "n", "", "Client application name.")
	cmd.Flags().StringSliceVarP(&imsConfig.RedirectURIs, "redirectURIs", "r", []string{}, "Redirect URIs (comma-separated or multiple flags).")

	return cmd
}
