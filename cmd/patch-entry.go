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
	"fmt"

	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// patchEntryCmd represents the entry command
var patchEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Partially updates an entry",
	Long: `Applies a partial update to an entry in Pleasant Password Server. Takes a JSON string as input.

To add a user access assignment to an entry, use --useraccess.
You can find available PermissionSetIds by running 'pleasant-cli get accesslevels'.

Examples:
pleasant-cli patch entry --id <id> --data '
{
    "Name": "NewNameForEntry"
}'

pleasant-cli patch entry --path 'Root/Folder1/TestEntry' --data '
{
    "Password": "MyNewPassword01"
}'

pleasant-cli patch entry --path 'Root/Folder1/TestEntry' --useraccess --data '
{
	"UserId": "788017e9-0ee0-460e-8de4-abb5016f65c5",
	"RoleId": "",
	"ZoneId": "",
	"PermissionSetId": "6fe3319c-21f0-48b0-a274-22fcca660de3",
	"AccessExpiry": "2020-12-31"
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

		var identifier string

		if cmd.Flags().Changed("path") {
			resourcePath, err := cmd.Flags().GetString("path")
			if err != nil {
				pleasant.ExitFatal(err)
			}

			id, err := pleasant.GetIdByResourcePath(baseUrl, resourcePath, "entry", bearerToken)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			identifier = id
		} else {
			id, err := cmd.Flags().GetString("id")
			if err != nil {
				pleasant.ExitFatal(err)
			}

			identifier = id
		}

		subPath := pleasant.PathEntry + "/" + identifier

		var msg string

		if cmd.Flags().Changed("useraccess") {
			subPath = subPath + "/useraccess"

			msg = fmt.Sprintf("User access assignment for entry %v added", identifier)

			_, err := pleasant.PostJsonString(baseUrl, subPath, json, bearerToken)
			if err != nil {
				pleasant.ExitFatal(err)
			}
		} else {
			msg = fmt.Sprintf("Existing entry with id %v patched", identifier)

			_, err = pleasant.PatchJsonString(baseUrl, subPath, json, bearerToken)
			if err != nil {
				pleasant.ExitFatal(err)
			}
		}

		pleasant.Exit(msg)
	},
}

func init() {
	patchCmd.AddCommand(patchEntryCmd)

	patchEntryCmd.Flags().StringP("path", "p", "", "Path to entry")
	patchEntryCmd.Flags().StringP("id", "i", "", "Id of entry")
	patchEntryCmd.MarkFlagsMutuallyExclusive("path", "id")
	patchEntryCmd.MarkFlagsOneRequired("path", "id")

	patchEntryCmd.Flags().StringP("data", "d", "", "JSON string with partial update/user access assignment")
	patchEntryCmd.MarkFlagRequired("data")

	patchEntryCmd.Flags().Bool("useraccess", false, "Add user access assignment to the entry")
}
