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

// getFoldersCmd represents the folders command
var getFoldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "Gets the entire folder tree",
	Long: `Gets the entire Pleasant Password tree.
Warning: this command can take a while to complete.
	
Example:
pleasant-cli get folders`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			return
		}

		baseUrl := viper.GetString("serverurl")
		bearerToken := viper.GetString("bearertoken.accesstoken")

		folder, err := pleasant.GetJsonBody(baseUrl, pleasant.PathFolders, bearerToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(folder)
	},
}

func init() {
	getCmd.AddCommand(getFoldersCmd)
}
