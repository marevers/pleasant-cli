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

// patchEntryCmd represents the entry command
var patchEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Partially updates an entry",
	Long: `Applies a partial update to an entry in Pleasant Password Server. Takes a JSON string as input.

Examples:
pleasant-cli patch entry --id <id> --data '
{
    "Name": "NewNameForEntry"
}'

pleasant-cli patch entry --path 'Root/Folder1/TestEntry' --data '
{
    "Password": "MyNewPassword01"
}'`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			return
		}

		baseUrl := viper.GetString("serverurl")
		bearerToken := viper.GetString("bearertoken.accesstoken")

		json, err := cmd.Flags().GetString("data")
		if err != nil {
			fmt.Println(err)
			return
		}

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

		subPath := pleasant.PathEntry + "/" + identifier

		_, err = pleasant.PatchJsonString(baseUrl, subPath, json, bearerToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Existing entry with id", identifier, "patched")
	},
}

func init() {
	patchCmd.AddCommand(patchEntryCmd)

	patchEntryCmd.Flags().StringP("path", "p", "", "Path to entry")
	patchEntryCmd.Flags().StringP("id", "i", "", "Id of entry")
	patchEntryCmd.MarkFlagsMutuallyExclusive("path", "id")
	patchEntryCmd.MarkFlagsOneRequired("path", "id")

	patchEntryCmd.Flags().StringP("data", "d", "", "JSON string with partial update")
	patchEntryCmd.MarkFlagRequired("data")
}
