package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets entries, folders, access levels, server info or password strength",
	Long:  `Gets entries, folders, access levels, server info or password strength`,
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().Bool("pretty", false, "Pretty-prints the JSON output")
}
