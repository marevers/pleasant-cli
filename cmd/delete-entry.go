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

// deleteEntryCmd represents the entry command
var deleteEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Archives or deletes an entry",
	Long: `Archives or deletes an entry from the Pleasant Password tree by its id or path.
A path must be absolute and starts with 'Root/', e.g. 'Root/Folder1/Folder2/Entry'.

By default, the entry is archived. If it should be deleted, use --delete.
WARNING: Deletion is permanent, use at your own risk.

Examples:
pleasant-cli delete entry --id <id>
pleasant-cli delete entry --path <path>
pleasant-cli delete entry --id <id> --delete`,
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

			id, err := pleasant.GetIdByResourcePath(baseUrl, resourcePath, "entry", bearerToken)
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

		subPath := pleasant.PathEntry + "/" + identifier

		_, err := pleasant.DeleteJsonString(baseUrl, subPath, json, bearerToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Entry with id", identifier, "archived/deleted")
	},
}

func init() {
	deleteCmd.AddCommand(deleteEntryCmd)

	deleteEntryCmd.Flags().StringP("path", "p", "", "Path to entry")
	deleteEntryCmd.Flags().StringP("id", "i", "", "Id of entry")
	deleteEntryCmd.MarkFlagsMutuallyExclusive("path", "id")
	deleteEntryCmd.MarkFlagsOneRequired("path", "id")

	deleteEntryCmd.Flags().Bool("delete", false, "Deletes the entry instead of archiving")
}
