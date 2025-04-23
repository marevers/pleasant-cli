package cmd

import (
	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// timeoutCmd represents the timeout command
var timeoutCmd = &cobra.Command{
	Use:   "timeout",
	Short: "Sets the Pleasant Password server timeout for pleasant-cli",
	Long: `Sets the Pleasant Password server timeout for pleasant-cli
It is specified as seconds and the default value (when it is unconfigured / set to 0 in the config file) is 20 seconds.

Example:
pleasant-cli config timeout 30`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		err := pleasant.WriteConfigFile(cfgFile, "Timeout", args[0])
		if err != nil {
			pleasant.ExitFatal(err)
		}

		pleasant.Exit("Timeout saved to:", cfgFile)
	},
}

func init() {
	configCmd.AddCommand(timeoutCmd)
}
