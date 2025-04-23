package cmd

import (
	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// serverurlCmd represents the serverurl command
var serverurlCmd = &cobra.Command{
	Use:   "serverurl",
	Short: "Sets the Pleasant Password server URL for pleasant-cli",
	Long: `Sets the Pleasant Password server URL for pleasant-cli
It must be specified as <PROTOCOL>://<URL>(:<PORT>).

If the port is either 80 or 443, it can be inferred from the protocol and can be omitted.

Example:
pleasant-cli config serverurl <SERVER URL>`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		err := pleasant.WriteConfigFile(cfgFile, "ServerUrl", args[0])
		if err != nil {
			pleasant.ExitFatal(err)
		}

		pleasant.Exit("Server URL saved to:", cfgFile)
	},
}

func init() {
	configCmd.AddCommand(serverurlCmd)
}
