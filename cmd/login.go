package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marevers/pleasant-cli/pleasant"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Pleasant Password Server",
	Long: `Log into Pleasant Password Server with username and password.
Username and password can either be entered interactively or by using flags.

Examples:
pleasant-cli login
pleasant-cli login --username <USERNAME> --password <PASSWORD>`,
	Run: func(cmd *cobra.Command, args []string) {
		if !pleasant.CheckPrerequisites(pleasant.IsServerUrlSet()) {
			pleasant.ExitFatal(pleasant.ErrPrereqNotMet)
		}

		var username string
		var password string

		if cmd.Flag("username").Value.String() != "" || cmd.Flag("password").Value.String() != "" {
			username = cmd.Flag("username").Value.String()
			password = cmd.Flag("password").Value.String()
		} else {
			username = pleasant.StringPrompt("Enter username:")
			password = pleasant.PasswordPrompt("Enter password:")
		}

		fmt.Println("Logging in to Pleasant Password Server...")

		bearerToken, err := pleasant.GetBearerToken(viper.GetString("serverUrl"), username, password)
		if errors.Is(err, pleasant.ErrBadRequest) {
			pleasant.ExitFatal(pleasant.ErrInvalidCredentials)
		} else if err != nil {
			pleasant.ExitFatal(err)
		}

		t := time.Now()
		ea := t.Unix() + int64(bearerToken.ExpiresIn)

		err = pleasant.WriteTokenFile(tokenFile, bearerToken.AccessToken, ea)
		if err != nil {
			pleasant.ExitFatal(err)
		}

		pleasant.Exit("Successfully logged in. Token saved to:", tokenFile, ", valid until", time.Unix(ea, 0))
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username for Pleasant Password Server")
	loginCmd.Flags().StringP("password", "p", "", "Password for Pleasant Password Server")
	loginCmd.MarkFlagsRequiredTogether("username", "password")
}
