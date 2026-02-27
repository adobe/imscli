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

type tokenDef struct {
	use        string
	alias      string
	label      string
	flagName   string
	field      *string
	successMsg string
	extraFlags func(cmd *cobra.Command, imsConfig *ims.Config)
}

func tokenCmd(imsConfig *ims.Config, def tokenDef) *cobra.Command {
	cmd := &cobra.Command{
		Use:     def.use,
		Aliases: []string{def.alias},
		Short:   fmt.Sprintf("Invalidate %s.", def.label),
		Long:    fmt.Sprintf("Invalidate %s.", def.label),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			err := imsConfig.InvalidateToken()
			if err != nil {
				return fmt.Errorf("error invalidating the %s: %w", def.label, err)
			}
			fmt.Println(def.successMsg)
			return nil
		},
	}

	cmd.Flags().StringVarP(def.field, def.flagName, "t", "", def.label+".")
	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS Client ID.")
	if def.extraFlags != nil {
		def.extraFlags(cmd, imsConfig)
	}

	return cmd
}

func AccessTokenCmd(imsConfig *ims.Config) *cobra.Command {
	return tokenCmd(imsConfig, tokenDef{
		use: "accessToken", alias: "acc", label: "access token",
		flagName: "accessToken", field: &imsConfig.AccessToken,
		successMsg: "Token invalidated successfully.",
	})
}

func RefreshTokenCmd(imsConfig *ims.Config) *cobra.Command {
	return tokenCmd(imsConfig, tokenDef{
		use: "refreshToken", alias: "ref", label: "refresh token",
		flagName: "refreshToken", field: &imsConfig.RefreshToken,
		successMsg: "Refresh token successfully invalidated.",
		extraFlags: func(cmd *cobra.Command, c *ims.Config) {
			cmd.Flags().BoolVarP(&c.Cascading, "cascading", "a", false,
				"Also invalidate all tokens obtained with the refresh token.")
		},
	})
}

func DeviceTokenCmd(imsConfig *ims.Config) *cobra.Command {
	return tokenCmd(imsConfig, tokenDef{
		use: "deviceToken", alias: "dev", label: "device token",
		flagName: "deviceToken", field: &imsConfig.DeviceToken,
		successMsg: "Token invalidated successfully.",
		extraFlags: func(cmd *cobra.Command, c *ims.Config) {
			cmd.Flags().BoolVarP(&c.Cascading, "cascading", "a", false,
				"Also invalidate all tokens obtained with the device token.")
		},
	})
}

func ServiceTokenCmd(imsConfig *ims.Config) *cobra.Command {
	return tokenCmd(imsConfig, tokenDef{
		use: "serviceToken", alias: "svc", label: "service token",
		flagName: "serviceToken", field: &imsConfig.ServiceToken,
		successMsg: "Service token successfully invalidated.",
		extraFlags: func(cmd *cobra.Command, c *ims.Config) {
			cmd.Flags().StringVarP(&c.ClientSecret, "clientSecret", "s", "", "IMS Client Secret.")
		},
	})
}
