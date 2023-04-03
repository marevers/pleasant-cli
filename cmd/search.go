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

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for entries and folders matching a query",
	Long: `Search for entries and folders matching a query.
	
Example: pleasant-cli search --query 'MyTestEntry'`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet, pleasant.IsTokenValid) {
			return
		}

		baseUrl := viper.GetString("serverurl")
		bearerToken := viper.GetString("bearertoken.accesstoken")

		query, err := cmd.Flags().GetString("query")
		if err != nil {
			fmt.Println(err)
			return
		}

		result, err := pleasant.PostSearch(baseUrl, query, bearerToken)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("query", "q", "", "Search query string")
	searchCmd.MarkFlagRequired("query")
}
