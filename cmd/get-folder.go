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

// getFolderCmd represents the folder command
var getFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Gets a folder and its entries by its id or path",
	Long: `Gets a folder and its entries by its id or path.
A path must be absolute and starts with 'Root/', e.g. 'Root/Folder1/Folder2/Folder3'.
	
Examples:
pleasant-cli get folder --id <id>
pleasant-cli get folder --path <path>`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			return
		}

		baseUrl := viper.GetString("serverurl")
		bearerToken := viper.GetString("bearertoken.accesstoken")

		var identifier string

		if cmd.Flags().Changed("path") {
			resourcePath, err := cmd.Flags().GetString("path")
			if err != nil {
				fmt.Println(err)
				return
			}

			id, err := pleasant.GetIdByResourcePath(baseUrl, resourcePath, "folder", bearerToken)
			if err != nil {
				fmt.Println(err)
				return
			}

			identifier = id
		} else {
			id, err := cmd.Flags().GetString("id")
			if err != nil {
				fmt.Println(err)
				return
			}

			identifier = id
		}

		subPath := pleasant.PathFolders + "/" + identifier

		if cmd.Flags().Changed("useraccess") {
			subPath = subPath + "/useraccess"
		}

		folder, err := pleasant.GetJsonBody(baseUrl, subPath, bearerToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		if cmd.Flags().Changed("pretty") {
			output, err := pleasant.PrettyPrintJson(folder)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(output)
			return
		}

		fmt.Println(folder)
	},
}

func init() {
	getCmd.AddCommand(getFolderCmd)

	getFolderCmd.Flags().StringP("path", "p", "", "Path to folder")
	getFolderCmd.Flags().StringP("id", "i", "", "Id of folder")
	getFolderCmd.MarkFlagsMutuallyExclusive("path", "id")
	getFolderCmd.MarkFlagsOneRequired("path", "id")

	getFolderCmd.Flags().Bool("useraccess", false, "Gets the users that have access to the folder")
}
