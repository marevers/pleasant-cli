package cmd

import (
	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// getServerInfoCmd represents the serverinfo command
var getServerInfoCmd = &cobra.Command{
	Use:   "serverinfo",
	Short: "Gets information about the server",
	Long: `Gets information about the Pleasant Password Server instance like the version,
OS information and DNS settings.

Example:
pleasant-cli get serverinfo`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			pleasant.ExitFatal(pleasant.ErrPrereqNotMet)
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

		serverInfo, err := pleasant.GetJsonBody(baseUrl, pleasant.PathServerInfo, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		if cmd.Flags().Changed("pretty") {
			output, err := pleasant.PrettyPrintJson(serverInfo)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pleasant.Exit(output)
		}

		pleasant.Exit(serverInfo)
	},
}

func init() {
	getCmd.AddCommand(getServerInfoCmd)
}
