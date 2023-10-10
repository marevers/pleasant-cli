/*
Copyright © 2023 Martijn Evers

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var tokenFile string

// Pleasant-CLI version
var version = "v0.0.3"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pleasant-cli",
	Version: version,
	Short:   "A simple CLI to interact with Pleasant Password Server.",
	Long: `pleasant-cli is an easy to use CLI that uses the Pleasant Password Server
API to interact with a Pleasant Password Server instance.

To use pleasant-cli, you must first set your server URL by running the following command:
pleasant-cli config serverurl <SERVER URL>

You can then log in by running:
pleasant-cli login`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pleasant-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token", "", "token file (default is $HOME/.pleasant-token.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pleasant-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pleasant-cli")

		// Set config file path
		cfgFile = filepath.Join(home, ".pleasant-cli.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	// Use token file from the flag
	if tokenFile != "" {
		// Use token file from the flag
		viper.AddConfigPath(tokenFile)
		viper.SetConfigType("ÿaml")
		viper.SetConfigName(filepath.Base(tokenFile))
		viper.MergeInConfig()
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search token file in home directory with name ".pleasant-token" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pleasant-token")
		viper.MergeInConfig()

		// Set token file path
		tokenFile = filepath.Join(home, ".pleasant-token.yaml")
	}
}
