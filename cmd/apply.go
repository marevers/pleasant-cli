package cmd

import (
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies a configuration to entries or folders",
	Long: `Applies a configuration to entries or folders.
Creates the entry or folder if it does not yet exist.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
