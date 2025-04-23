package cmd

import (
	"os"

	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// clearTokenCmd represents the cleartoken command
var clearTokenCmd = &cobra.Command{
	Use:   "cleartoken",
	Short: "Removes the token file",
	Long: `Removes the token file
	
Example:
pleasant-cli config cleartoken
pleasant-cli config cleartoken --token <PATH>`,
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Remove(tokenFile)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		pleasant.Exit("Token file deleted:", tokenFile)
	},
}

func init() {
	configCmd.AddCommand(clearTokenCmd)
}
