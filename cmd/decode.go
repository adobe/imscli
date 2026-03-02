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
	"os"
	"time"

	"github.com/adobe/imscli/cmd/pretty"
	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func decodeCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "decode",
		Aliases: []string{"dec"},
		Short:   "Decode a JWT token.",
		Long:    "Decode a JWT token and display the header and payload as prettified JSON.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			decoded, err := imsConfig.DecodeToken()
			if err != nil {
				return fmt.Errorf("error decoding the token: %w", err)
			}

			output := fmt.Sprintf(`{"header":%s,"payload":%s}`, decoded.Header, decoded.Payload)
			fmt.Println(pretty.JSON(output))

			// When verbose, show human-readable token expiration on stderr
			// so it doesn't pollute the JSON output on stdout.
			if imsConfig.Verbose {
				printTokenExpiration(decoded.Payload)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.Token, "token", "t", "", "Token.")

	return cmd
}

// printTokenExpiration parses the "exp" claim from a JWT payload and prints
// a human-readable expiration message to stderr.
func printTokenExpiration(payload string) {
	var claims map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &claims); err != nil {
		return
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return
	}

	expTime := time.Unix(int64(exp), 0).UTC()
	now := time.Now().UTC()

	if now.After(expTime) {
		fmt.Fprintf(os.Stderr, "\nToken expired: %s (%s ago)\n",
			expTime.Format(time.RFC3339), now.Sub(expTime).Truncate(time.Second))
	} else {
		fmt.Fprintf(os.Stderr, "\nToken expires: %s (in %s)\n",
			expTime.Format(time.RFC3339), expTime.Sub(now).Truncate(time.Second))
	}
}
