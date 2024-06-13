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
	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// getAccesslevelsCmd represents the accesslevels command
var getAccesslevelsCmd = &cobra.Command{
	Use:   "accesslevels",
	Short: "Gets a list of access levels",
	Long: `Gets a list of access levels from Pleasant Password Server.
	
Example:
pleasant-cli get accesslevels`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			pleasant.ExitFatal(pleasant.ErrPrereqNotMet)
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

		accesslevels, err := pleasant.GetJsonBody(baseUrl, pleasant.PathAccessLevels, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		if cmd.Flags().Changed("pretty") {
			output, err := pleasant.PrettyPrintJson(accesslevels)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pleasant.Exit(output)
		}

		pleasant.Exit(accesslevels)
	},
}

func init() {
	getCmd.AddCommand(getAccesslevelsCmd)
}
