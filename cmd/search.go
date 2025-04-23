package cmd

import (
	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for entries and folders matching a query",
	Long: `Search for entries and folders matching a query.
	
Example:
pleasant-cli search --query 'MyTestEntry'`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			pleasant.ExitFatal(pleasant.ErrPrereqNotMet)
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

		query, err := cmd.Flags().GetString("query")
		if err != nil {
			pleasant.ExitFatal(err)
		}

		result, err := pleasant.PostSearch(baseUrl, query, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		if cmd.Flags().Changed("pretty") {
			output, err := pleasant.PrettyPrintJson(result)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pleasant.Exit(output)
		}

		pleasant.Exit(result)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("query", "q", "", "Search query string")
	searchCmd.MarkFlagRequired("query")

	searchCmd.Flags().Bool("pretty", false, "Pretty-prints the JSON output")
}
