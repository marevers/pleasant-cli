package cmd

import (
	"github.com/spf13/cobra"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Partially updates entries or folders or adds user access assignments for them",
	Long:  `Partially updates entries or folders or adds user access assignments for them`,
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(patchCmd)
}
