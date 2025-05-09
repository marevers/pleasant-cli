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

To get the username of an entry, use --username.
To get the password of an entry, use --password.
For these two options, if you include the flag --clip, the output
will be copied to your clipboard instead.

To get the attachments of an entry, use --attachments.

Examples:
pleasant-cli get entry --id <id>
pleasant-cli get entry --path <path>
pleasant-cli get entry --id <id> --username
pleasant-cli get entry --id <id> --password --clip
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

			if cmd.Flags().Changed("clip") {
				err := pleasant.CopyToClipboard(en.Username)
				if err != nil {
					pleasant.ExitFatal(err)
				}

				pleasant.Exit("Username copied to clipboard")
			}

			pleasant.Exit(en.Username)
		case cmd.Flags().Changed("password"):
			if cmd.Flags().Changed("clip") {
				err := pleasant.CopyToClipboard(pleasant.Unescape(pleasant.TrimDoubleQuotes(entry)))
				if err != nil {
					pleasant.ExitFatal(err)
				}

				pleasant.Exit("Password copied to clipboard")
			}

			pleasant.Exit(pleasant.Unescape(pleasant.TrimDoubleQuotes(entry)))
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

	getEntryCmd.RegisterFlagCompletionFunc("path", func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return pleasant.CompletePathFlag(toComplete, true)
	})

	getEntryCmd.Flags().Bool("username", false, "Get the username of the entry")
	getEntryCmd.Flags().Bool("password", false, "Get the password of the entry")
	getEntryCmd.Flags().Bool("attachments", false, "Gets the attachments of the entry")
	getEntryCmd.Flags().Bool("useraccess", false, "Gets the users that have access to the entry")
	getEntryCmd.MarkFlagsMutuallyExclusive("username", "password", "attachments", "useraccess")

	getEntryCmd.Flags().Bool("clip", false, "Copy the output to the clipboard instead (username or password)")
}
