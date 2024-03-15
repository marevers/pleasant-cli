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
)

// getRootfolderCmd represents the rootfolder command
var getRootfolderCmd = &cobra.Command{
	Use:   "rootfolder",
	Short: "Gets the id of the root folder",
	Long: `Returns the id of the root folder in the Pleasant Password Server tree

Example:
pleasant-cli get rootfolder`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			return
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

		rootFolderId, err := pleasant.GetJsonBody(baseUrl, pleasant.PathRootFolder, bearerToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		if cmd.Flags().Changed("pretty") {
			output, err := pleasant.PrettyPrintJson(rootFolderId)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(output)
			return
		}

		fmt.Println(rootFolderId)
	},
}

func init() {
	getCmd.AddCommand(getRootfolderCmd)
}
