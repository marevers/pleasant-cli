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
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marevers/pleasant-cli/pleasant"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Pleasant Password Server",
	Long: `Log into Pleasant Password Server with username and password.
Username and password can either be entered interactively or by using flags.

Examples:
pleasant-cli login
pleasant-cli login --username <USERNAME> --password <PASSWORD>`,
	Run: func(cmd *cobra.Command, args []string) {
		var username string
		var password string

		if cmd.Flag("username").Value.String() != "" || cmd.Flag("password").Value.String() != "" {
			username = cmd.Flag("username").Value.String()
			password = cmd.Flag("password").Value.String()
		} else {
			username = pleasant.StringPrompt("Enter username:")
			password = pleasant.PasswordPrompt("Enter password:")
		}

		fmt.Println("\nLogging in to Pleasant Password Server...")

		bearerToken, err := pleasant.GetBearerToken(viper.GetString("serverUrl"), username, password)
		if errors.Is(err, pleasant.ErrBadRequest) {
			fmt.Println(pleasant.ErrInvalidCredentials)
			return
		} else if err != nil {
			fmt.Println(err)
			return
		}

		t := time.Now()
		ea := t.Unix() + int64(bearerToken.ExpiresIn)

		err = pleasant.WriteTokenFile(tokenFile, bearerToken.AccessToken, ea)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Successfully logged in. Token saved to:", tokenFile, ", valid until", time.Unix(ea, 0))
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username for Pleasant Password Server")
	loginCmd.Flags().StringP("password", "p", "", "Password for Pleasant Password Server")
	loginCmd.MarkFlagsRequiredTogether("username", "password")
}
