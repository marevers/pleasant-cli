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

// getEntryCmd represents the entry command
var getEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Gets an entry by its id or path",
	Long: `Gets an entry from the Pleasant Password tree by its id or path.
A path must be absolute and starts with 'Root/', e.g. 'Root/Folder1/Folder2/Entry'.

To get the password of an entry, use --password.
To get the attachments of an entry, use --attachments.

Examples:
pleasant-cli get entry --id <id>
pleasant-cli get entry --path <path>
pleasant-cli get entry --id <id> --password
pleasant-cli get entry --path <path> --attachments`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			pleasant.ExitFatal(pleasant.ErrPrereqNotMet)
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

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

		switch {
		case cmd.Flags().Changed("password"):
			subPath = subPath + "/password"
		case cmd.Flags().Changed("attachments"):
			subPath = subPath + "/attachments"
		case cmd.Flags().Changed("useraccess"):
			subPath = subPath + "/useraccess"
		}

		entry, err := pleasant.GetJsonBody(baseUrl, subPath, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		switch {
		case cmd.Flags().Changed("pretty"):
			output, err := pleasant.PrettyPrintJson(entry)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pleasant.Exit(output)
		case cmd.Flags().Changed("username"):
			en, err := pleasant.UnmarshalEntry(entry)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pleasant.Exit(en.Username)
		case cmd.Flags().Changed("password"):
			pleasant.Exit(pleasant.TrimDoubleQuotes(entry))
		default:
			pleasant.Exit(entry)
		}
	},
}

func init() {
	getCmd.AddCommand(getEntryCmd)

	getEntryCmd.Flags().StringP("path", "p", "", "Path to entry")
	getEntryCmd.Flags().StringP("id", "i", "", "Id of entry")
	getEntryCmd.MarkFlagsMutuallyExclusive("path", "id")
	getEntryCmd.MarkFlagsOneRequired("path", "id")

	getEntryCmd.Flags().Bool("username", false, "Get the username of the entry")
	getEntryCmd.Flags().Bool("password", false, "Get the password of the entry")
	getEntryCmd.Flags().Bool("attachments", false, "Gets the attachments of the entry")
	getEntryCmd.Flags().Bool("useraccess", false, "Gets the users that have access to the entry")
	getEntryCmd.MarkFlagsMutuallyExclusive("username", "password", "attachments", "useraccess")
}
