package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the create command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Archives or deletes entries or folders or user access assignments for them",
	Long:  `Archives or deletes entries or folders or user access assignments for them`,
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
