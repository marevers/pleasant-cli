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

// createFolderCmd represents the folder command
var createFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Creates a folder",
	Long: `Creates a folder in Pleasant Password Server. Takes a JSON string as input.
Returns its id if succesful.

Example: pleasant-cli create folder --data '
{
    "CustomUserFields": {},
    "CustomApplicationFields": {},
    "Children": [],
    "Credentials": [],
    "Tags": [],
    "Name": "TestFolder",
    "ParentId": "c04f874b-90f7-4b33-97d0-a92e011fb712",
    "Notes": null,
    "Expires": null
}'`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet, pleasant.IsTokenValid) {
			return
		}

		baseUrl := viper.GetString("serverurl")
		bearerToken := viper.GetString("bearertoken.accesstoken")

		json, err := cmd.Flags().GetString("data")
		if err != nil {
			fmt.Println(err)
			return
		}

		id, err := pleasant.PostJsonString(baseUrl, pleasant.PathFolders, json, bearerToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(id)
	},
}

func init() {
	createCmd.AddCommand(createFolderCmd)

	createFolderCmd.Flags().StringP("data", "d", "", "JSON string with folder data")
	createFolderCmd.MarkFlagRequired("data")
}
