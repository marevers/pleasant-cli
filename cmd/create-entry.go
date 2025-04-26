package cmd

import (
	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// createEntryCmd represents the entry command
var createEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Creates an entry",
	Long: `Creates an entry in Pleasant Password Server. Takes a JSON string as input.
Returns its id if succesful.
'GroupId' can be omitted if the path of the entry is supplied.

Examples:
pleasant-cli create entry --data '
{
    "CustomUserFields": {},
    "CustomApplicationFields": {},
    "Tags": [],
    "Name": "TestEntry",
    "Username": "MyUserName",
    "Password": "MyPassword01",
    "Url": "",
    "Notes": "",
    "GroupId": "c04f874b-90f7-4b33-97d0-a92e011fb712",
    "Expires": null
}'

pleasant-cli create entry --path 'Root/Folder1/TestEntry' --data '
{
    "CustomUserFields": {},
    "CustomApplicationFields": {},
    "Tags": [],
    "Name": "TestEntry",
    "Username": "MyUserName",
    "Password": "MyPassword01",
    "Url": "",
    "Notes": "",
    "Expires": null
}'`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			pleasant.ExitFatal(pleasant.ErrPrereqNotMet)
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

		json, err := cmd.Flags().GetString("data")
		if err != nil {
			pleasant.ExitFatal(err)
		}

		if cmd.Flags().Changed("path") {
			resourcePath, err := cmd.Flags().GetString("path")
			if err != nil {
				pleasant.ExitFatal(err)
			}

			input, err := pleasant.UnmarshalEntry(json)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			pid, err := pleasant.GetParentIdByResourcePath(baseUrl, resourcePath, bearerToken)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			if !pleasant.PathAndNameMatching(resourcePath, input.Name) {
				pleasant.ExitFatal("error: entry name from path and data do not match")
			}

			input.GroupId = pid

			j, err := pleasant.MarshalEntry(input)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			json = j
		}

		if cmd.Flags().Changed("no-duplicates") {
			exists, err := pleasant.DuplicateEntryExists(baseUrl, json, bearerToken)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			if exists {
				pleasant.ExitFatal(pleasant.ErrDuplicateEntry)
			}
		}

		id, err := pleasant.PostJsonString(baseUrl, pleasant.PathEntry, json, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		pleasant.Exit(id)
	},
}

func init() {
	createCmd.AddCommand(createEntryCmd)

	createEntryCmd.Flags().StringP("data", "d", "", "JSON string with entry data")
	createEntryCmd.Flags().StringP("path", "p", "", "Path to entry")
	createEntryCmd.MarkFlagRequired("data")

	createEntryCmd.RegisterFlagCompletionFunc("path", func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return pleasant.CompletePathFlag(toComplete, false)
	})

	createEntryCmd.Flags().Bool("no-duplicates", false, "Avoid creating duplicate entries")
}
