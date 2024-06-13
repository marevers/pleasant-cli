/*
Copyright © 2024 Martijn Evers

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

// applyEntryCmd represents the entry command
var applyEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Applies a configuration to an entry",
	Long: `Applies a configuration to an entry in Pleasant Password Server. Takes a JSON string as input.
If the entry does not exist, it is created using the supplied configuration
and the command returns the entry's id. If it exists, any changes will be applied to the entry.
Partial updates are allowed e.g. only supplying 'Name' and 'Password' to update the password.
Note: you cannot use this command to change the name of an entry.

'GroupId' can be omitted if the path of the entry is supplied.

Examples:
pleasant-cli apply entry --data '
{
    "Name": "TestEntry",
    "Username": "MyNewUserName",
    "Password": "MyNewPassword01",
    "GroupId": "c04f874b-90f7-4b33-97d0-a92e011fb712"
}'

pleasant-cli apply entry --path 'Root/Folder1/TestEntry' --data '
{
    "Name": "TestEntry",
    "Username": "MyNewUserName",
    "Password": "MyNewPassword01"
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

		input, err := pleasant.UnmarshalEntry(json)
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
				pleasant.ExitFatal("error: entry name from path and data do not match")
			}

			input.GroupId = pid

			j, err := pleasant.MarshalEntry(input)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			json = j
		}

		id, err := pleasant.DuplicateEntryId(baseUrl, json, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		if id != "" {
			subPath := pleasant.PathEntry + "/" + id

			_, err := pleasant.PatchJsonString(baseUrl, subPath, json, bearerToken)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pleasant.Exit("Existing entry with id", id, "patched")
		}

		id, err = pleasant.PostJsonString(baseUrl, pleasant.PathEntry, json, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		pleasant.Exit(id)
	},
}

func init() {
	applyCmd.AddCommand(applyEntryCmd)

	applyEntryCmd.Flags().StringP("data", "d", "", "JSON string with entry data")
	applyEntryCmd.Flags().StringP("path", "p", "", "Path to entry")
	applyEntryCmd.MarkFlagRequired("data")
}
