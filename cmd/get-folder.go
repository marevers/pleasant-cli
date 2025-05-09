package cmd

import (
	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// getFolderCmd represents the folder command
var getFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Gets a folder and its entries by its id or path",
	Long: `Gets a folder and its entries by its id or path.
A path must be absolute and starts with 'Root/', e.g. 'Root/Folder1/Folder2/Folder3'.
	
Examples:
pleasant-cli get folder --id <id>
pleasant-cli get folder --path <path>`,
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

			id, err := pleasant.GetIdByResourcePath(baseUrl, resourcePath, "folder", bearerToken)
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

		subPath := pleasant.PathFolders + "/" + identifier

		if cmd.Flags().Changed("useraccess") {
			subPath = subPath + "/useraccess"
		}

		folder, err := pleasant.GetJsonBody(baseUrl, subPath, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		if cmd.Flags().Changed("pretty") {
			output, err := pleasant.PrettyPrintJson(folder)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pleasant.Exit(output)
		}

		pleasant.Exit(folder)
	},
}

func init() {
	getCmd.AddCommand(getFolderCmd)

	getFolderCmd.Flags().StringP("path", "p", "", "Path to folder")
	getFolderCmd.Flags().StringP("id", "i", "", "Id of folder")
	getFolderCmd.MarkFlagsMutuallyExclusive("path", "id")
	getFolderCmd.MarkFlagsOneRequired("path", "id")

	getFolderCmd.RegisterFlagCompletionFunc("path", func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return pleasant.CompletePathFlag(toComplete, false)
	})

	getFolderCmd.Flags().Bool("useraccess", false, "Gets the users that have access to the folder")
}
