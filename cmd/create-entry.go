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

// createEntryCmd represents the entry command
var createEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Creates an entry",
	Long: `Creates an entry in Pleasant Password Server. Takes a JSON string as input.
Returns its id if succesful.
'GroupId' can be omitted if the path of the entry is supplied.

Examples:
pleasant-cli create entry --data '
{
    "CustomUserFields": {},
    "CustomApplicationFields": {},
    "Tags": [],
    "Name": "TestEntry",
    "Username": "MyUserName",
    "Password": "MyPassword01",
    "Url": "",
    "Notes": "",
    "GroupId": "c04f874b-90f7-4b33-97d0-a92e011fb712",
    "Expires": null
}'

pleasant-cli create entry --path 'Root/Folder1/TestEntry' --data '
{
    "CustomUserFields": {},
    "CustomApplicationFields": {},
    "Tags": [],
    "Name": "TestEntry",
    "Username": "MyUserName",
    "Password": "MyPassword01",
    "Url": "",
    "Notes": "",
    "Expires": null
}'`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			return
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

		json, err := cmd.Flags().GetString("data")
		if err != nil {
			fmt.Println(err)
			return
		}

		if cmd.Flags().Changed("path") {
			resourcePath, err := cmd.Flags().GetString("path")
			if err != nil {
				fmt.Println(err)
				return
			}

			input, err := pleasant.UnmarshalEntry(json)
			if err != nil {
				fmt.Println(err)
				return
			}

			pid, err := pleasant.GetParentIdByResourcePath(baseUrl, resourcePath, bearerToken)
			if err != nil {
				fmt.Println(err)
				return
			}

			if !pleasant.PathAndNameMatching(resourcePath, input.Name) {
				fmt.Println("error: entry name from path and data do not match")
				return
			}

			input.GroupId = pid

			j, err := pleasant.MarshalEntry(input)
			if err != nil {
				fmt.Println(err)
				return
			}

			json = j
		}

		if cmd.Flags().Changed("no-duplicates") {
			exists, err := pleasant.DuplicateEntryExists(baseUrl, json, bearerToken)
			if err != nil {
				fmt.Println(err)
				return
			}

			if exists {
				fmt.Println(pleasant.ErrDuplicateEntry)
				return
			}
		}

		id, err := pleasant.PostJsonString(baseUrl, pleasant.PathEntry, json, bearerToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(id)
	},
}

func init() {
	createCmd.AddCommand(createEntryCmd)

	createEntryCmd.Flags().StringP("data", "d", "", "JSON string with entry data")
	createEntryCmd.Flags().StringP("path", "p", "", "Path to entry")
	createEntryCmd.MarkFlagRequired("data")

	createEntryCmd.Flags().Bool("no-duplicates", false, "Avoid creating duplicate entries")
}
