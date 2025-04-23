package cmd

import (
	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// getFoldersCmd represents the folders command
var getFoldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "Gets the entire folder tree",
	Long: `Gets the entire Pleasant Password tree.
WARNING: this command can take a while to complete.
	
Example:
pleasant-cli get folders`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			pleasant.ExitFatal(pleasant.ErrPrereqNotMet)
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

		folder, err := pleasant.GetJsonBody(baseUrl, pleasant.PathFolders, bearerToken)
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
	getCmd.AddCommand(getFoldersCmd)
}
