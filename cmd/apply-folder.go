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

// applyFolderCmd represents the folder command
var applyFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Applies a configuration to a folder",
	Long: `Applies a configuration to a folder in Pleasant Password Server. Takes a JSON string as input.
If the folder does not exist, it is created using the supplied configuration
and the command returns the folder's id. If it exists, any changes will be applied to the folder.
Partial updates are allowed e.g. only supplying 'Name' and 'Notes' to update the notes.
Note: you cannot use this command to change the name of a folder.

'ParentId' can be omitted if the path of the folder is supplied.

Examples:
pleasant-cli apply folder --data '
{
    "Notes": "New note for the folder",
    "Name": "TestFolder",
    "ParentId": "c04f874b-90f7-4b33-97d0-a92e011fb712"
}'

pleasant-cli apply folder --path 'Root/Folder1/TestFolder' --data '
{
    "Notes": "New note for the folder",
    "Name": "TestFolder"
}'`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			pleasant.ExitFatal(pleasant.ErrPrereqNotMet)
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

		json, err := cmd.Flags().GetString("data")
		if err != nil {
			pleasant.ExitFatal(err)
		}

		input, err := pleasant.UnmarshalFolder(json)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		if cmd.Flags().Changed("path") {
			resourcePath, err := cmd.Flags().GetString("path")
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pid, err := pleasant.GetParentIdByResourcePath(baseUrl, resourcePath, bearerToken)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			if !pleasant.PathAndNameMatching(resourcePath, input.Name) {
				pleasant.ExitFatal("error: folder name from path and data do not match")
			}

			input.ParentId = pid

			j, err := pleasant.MarshalFolder(input)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			json = j
		}

		id, err := pleasant.DuplicateFolderId(baseUrl, json, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		if id != "" {
			subPath := pleasant.PathFolders + "/" + id

			_, err := pleasant.PatchJsonString(baseUrl, subPath, json, bearerToken)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pleasant.Exit("Existing folder with id", id, "patched")
		}

		id, err = pleasant.PostJsonString(baseUrl, pleasant.PathFolders, json, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		pleasant.Exit(id)
	},
}

func init() {
	applyCmd.AddCommand(applyFolderCmd)

	applyFolderCmd.Flags().StringP("data", "d", "", "JSON string with folder data")
	applyFolderCmd.Flags().StringP("path", "p", "", "Path to folder")
	applyFolderCmd.MarkFlagRequired("data")
}
