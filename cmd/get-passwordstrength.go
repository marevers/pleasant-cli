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

	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getPasswordStrengthCmd represents the passwordstrength command
var getPasswordStrengthCmd = &cobra.Command{
	Use:   "passwordstrength",
	Short: "Gets the password strength for a password",
	Long: `Gets a numerical and descriptive representation of the strength of a password.

Example:
pleasant-cli get passwordstrength --password <PASSWORD>
pleasant-cli get passwordstrength -p <PASSWORD>`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			return
		}

		baseUrl := viper.GetString("serverurl")
		bearerToken := viper.GetString("bearertoken.accesstoken")

		pw, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println(err)
			return
		}

		json := fmt.Sprintf(`{"Password":"%v"}`, pw)

		pwStr, err := pleasant.PostJsonString(baseUrl, pleasant.PathPwStr, json, bearerToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		if cmd.Flags().Changed("pretty") {
			output, err := pleasant.PrettyPrintJson(pwStr)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(output)
			return
		}

		fmt.Println(pwStr)
	},
}

func init() {
	getCmd.AddCommand(getPasswordStrengthCmd)

	getPasswordStrengthCmd.Flags().StringP("password", "p", "", "Password to test password strength")
	getPasswordStrengthCmd.MarkFlagRequired("password")
}
