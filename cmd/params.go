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
	"os"

	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Centralize in one place the processing of:
//   - Environment variables.
//   - Configuration file.
//   - Command parameters.
//
// The processing function only needs to be run as the PersistentPreRunE in the root command.

func initParams(cmd *cobra.Command, params *ims.Config, configFile string) error {
	v := viper.New()

	// Setup env vars
	v.SetEnvPrefix("ims")
	v.AutomaticEnv()

	// Command flags (local + inherited persistent flags)
	err := v.BindPFlags(cmd.Flags())
	if err != nil {
		return fmt.Errorf("unable to process command flags: %w", err)
	}
	err = v.BindPFlags(cmd.InheritedFlags())
	if err != nil {
		return fmt.Errorf("unable to process inherited flags: %w", err)
	}

	if configFile == "" {
		// Configuration file ( ~/.config/imscli.ext )
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("unable to find configuration directory: %w", err)
		}

		v.AddConfigPath(".")
		v.AddConfigPath(configDir)
		v.SetConfigName("imscli")
		err = v.ReadInConfig()
		if err != nil {
			// Ignore ConfigFileNotFoundError, since config file is not mandatory
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return fmt.Errorf("unable to read configuration file: %w", err)
			}
		}
	} else {
		v.SetConfigFile(configFile)
		err = v.ReadInConfig()
		if err != nil {
			return fmt.Errorf("unable to read configuration file: %w", err)
		}
	}

	err = v.Unmarshal(params)
	if err != nil {
		return fmt.Errorf("unable to parse configuration file: %w", err)
	}

	return nil
}
