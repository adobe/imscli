// Copyright 2021 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

// Package validate implements the validate subcommands for token validation.
package validate

import (
	"fmt"

	"github.com/adobe/imscli/cmd/prettify"
	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

type tokenDef struct {
	use      string
	alias    string
	label    string
	flagName string
	field    *string
}

func tokenCmd(imsConfig *ims.Config, def tokenDef) *cobra.Command {
	cmd := &cobra.Command{
		Use:     def.use,
		Aliases: []string{def.alias},
		Short:   fmt.Sprintf("Validate %s.", def.label),
		Long:    fmt.Sprintf("Validate %s.", def.label),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			resp, err := imsConfig.ValidateToken()
			if err != nil {
				return fmt.Errorf("error validating the %s: %w", def.label, err)
			}
			if !resp.Valid {
				return fmt.Errorf("invalid token: %v", resp.Info)
			}
			fmt.Println(prettify.JSON(resp.Info))
			return nil
		},
	}

	cmd.Flags().StringVarP(def.field, def.flagName, "t", "", def.label+".")
	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS Client ID.")

	return cmd
}

func AccessTokenCmd(imsConfig *ims.Config) *cobra.Command {
	return tokenCmd(imsConfig, tokenDef{
		use: "accessToken", alias: "acc", label: "access token",
		flagName: "accessToken", field: &imsConfig.AccessToken,
	})
}

func RefreshTokenCmd(imsConfig *ims.Config) *cobra.Command {
	return tokenCmd(imsConfig, tokenDef{
		use: "refreshToken", alias: "ref", label: "refresh token",
		flagName: "refreshToken", field: &imsConfig.RefreshToken,
	})
}

func DeviceTokenCmd(imsConfig *ims.Config) *cobra.Command {
	return tokenCmd(imsConfig, tokenDef{
		use: "deviceToken", alias: "dev", label: "device token",
		flagName: "deviceToken", field: &imsConfig.DeviceToken,
	})
}

func AuthzCodeCmd(imsConfig *ims.Config) *cobra.Command {
	return tokenCmd(imsConfig, tokenDef{
		use: "authorizationCode", alias: "authz", label: "authorization code",
		flagName: "authorizationCode", field: &imsConfig.AuthorizationCode,
	})
}
