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
	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func authzCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "authorize",
		Aliases: []string{"authz"},
		Short:   "Negotiate an access token with IMS.",
		Long: `The authorize command negotiates an access token with IMS, performing one of the supported types of authorization.

This command has no effect by itself, the authorization type needs to be specified as a subcommand.
`,
	}
	cmd.AddCommand(
		authzUserCmd(imsConfig),
		authzServiceCmd(imsConfig),
		authzJWTCmd(imsConfig),
		authzClientCredentialsCmd(imsConfig),
	)
	return cmd
}
