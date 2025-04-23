package cmd

import (
	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// getRootfolderCmd represents the rootfolder command
var getRootfolderCmd = &cobra.Command{
	Use:   "rootfolder",
	Short: "Gets the id of the root folder",
	Long: `Returns the id of the root folder in the Pleasant Password Server tree

Example:
pleasant-cli get rootfolder`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			pleasant.ExitFatal(pleasant.ErrPrereqNotMet)
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

		rootFolderId, err := pleasant.GetJsonBody(baseUrl, pleasant.PathRootFolder, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		if cmd.Flags().Changed("pretty") {
			output, err := pleasant.PrettyPrintJson(rootFolderId)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pleasant.Exit(output)
		}

		pleasant.Exit(rootFolderId)
	},
}

func init() {
	getCmd.AddCommand(getRootfolderCmd)
}
