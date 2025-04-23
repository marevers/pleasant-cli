package cmd

import (
	"fmt"

	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// deleteFolderCmd represents the entry command
var deleteFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Archives or deletes a folder or a user access assignment for it",
	Long: `Archives or deletes a folder from the Pleasant Password tree by its id or path.
Anything contained in the folder is also archived/deleted.
A path must be absolute and starts with 'Root/', e.g. 'Root/Folder1/Folder2/Folder3'.
Instead of the folder, a user access assignment can also be archived or deleted by appending --useraccess <accessrowid>.

By default, the folder is archived. If it should be deleted, use --delete.
WARNING: Deletion is permanent, use at your own risk.

Examples:
pleasant-cli delete folder --id <id>
pleasant-cli delete folder --path <path>
pleasant-cli delete folder --id <id> --delete
pleasant-cli delete folder --id <id> --delete --useraccess <accessrowid>`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet(), pleasant.IsTokenValid()) {
			pleasant.ExitFatal(pleasant.ErrPrereqNotMet)
		}

		baseUrl, bearerToken := pleasant.LoadConfig()

		var identifier string

		if cmd.Flags().Changed("path") {
			resourcePath, err := cmd.Flags().GetString("path")
			if err != nil {
				pleasant.ExitFatal(err)
			}

			id, err := pleasant.GetIdByResourcePath(baseUrl, resourcePath, "folder", bearerToken)
			if err != nil {
				pleasant.ExitFatal(err)
			}

			identifier = id
		} else {
			id, err := cmd.Flags().GetString("id")
			if err != nil {
				pleasant.ExitFatal(err)
			}

			identifier = id
		}

		var action string

		if cmd.Flags().Changed("delete") {
			action = "Delete"
		} else {
			action = "Archive"
		}

		json := fmt.Sprintf(`{"Action":"%v","Comment":"Archived/deleted by Pleasant-CLI"}`, action)

		subPath := pleasant.PathFolders + "/" + identifier

		var msg string

		if cmd.Flags().Changed("useraccess") {
			ua, err := cmd.Flags().GetString("useraccess")
			if err != nil {
				pleasant.ExitFatal(err)
			}

			subPath = subPath + "/useraccess/" + ua

			msg = fmt.Sprintf("User access assignment %v deleted from folder %v", ua, identifier)
		} else {
			msg = fmt.Sprintf("Folder with id %v archived/deleted", identifier)
		}

		_, err := pleasant.DeleteJsonString(baseUrl, subPath, json, bearerToken)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		pleasant.Exit(msg)
	},
}

func init() {
	deleteCmd.AddCommand(deleteFolderCmd)

	deleteFolderCmd.Flags().StringP("path", "p", "", "Path to folder")
	deleteFolderCmd.Flags().StringP("id", "i", "", "Id of folder")
	deleteFolderCmd.MarkFlagsMutuallyExclusive("path", "id")
	deleteFolderCmd.MarkFlagsOneRequired("path", "id")

	deleteFolderCmd.Flags().String("useraccess", "", "Archives/deletes the user access assignment with this id")
	deleteFolderCmd.Flags().Bool("delete", false, "Deletes the folder instead of archiving")
}
