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

// deleteFolderCmd represents the entry command
var deleteFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Archives or deletes a folder or a user access assignment for it",
	Long: `Archives or deletes a folder from the Pleasant Password tree by its id or path.
Anything contained in the folder is also archived/deleted.
A path must be absolute and starts with 'Root/', e.g. 'Root/Folder1/Folder2/Folder3'.
Instead of the folder, a user access assignment can also be archived or deleted by appending --useraccess <accessrowid>.

By default, the folder is archived. If it should be deleted, use --delete.
WARNING: Deletion is permanent, use at your own risk.

Examples:
pleasant-cli delete folder --id <id>
pleasant-cli delete folder --path <path>
pleasant-cli delete folder --id <id> --delete
pleasant-cli delete folder --id <id> --delete --useraccess <accessrowid>`,
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

		var action string

		if cmd.Flags().Changed("delete") {
			action = "Delete"
		} else {
			action = "Archive"
		}

		json := fmt.Sprintf(`{"Action":"%v","Comment":"Archived/deleted by Pleasant-CLI"}`, action)

		subPath := pleasant.PathFolders + "/" + identifier

		var msg string

		if cmd.Flags().Changed("useraccess") {
			ua, err := cmd.Flags().GetString("useraccess")
			if err != nil {
				fmt.Println(err)
				return
			}

			subPath = subPath + "/useraccess/" + ua

			msg = fmt.Sprintf("User access assignment %v deleted from folder %v", ua, identifier)
		} else {
			msg = fmt.Sprintf("Folder with id %v archived/deleted", identifier)
		}

		_, err := pleasant.DeleteJsonString(baseUrl, subPath, json, bearerToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(msg)
	},
}

func init() {
	deleteCmd.AddCommand(deleteFolderCmd)

	deleteFolderCmd.Flags().StringP("path", "p", "", "Path to folder")
	deleteFolderCmd.Flags().StringP("id", "i", "", "Id of folder")
	deleteFolderCmd.MarkFlagsMutuallyExclusive("path", "id")
	deleteFolderCmd.MarkFlagsOneRequired("path", "id")

	deleteFolderCmd.Flags().String("useraccess", "", "Archives/deletes the user access assignment with this id")
	deleteFolderCmd.Flags().Bool("delete", false, "Deletes the folder instead of archiving")
}
