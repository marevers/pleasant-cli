/*
Copyright © 2023 Martijn Evers

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/marevers/pleasant-cli/pleasant"
	"github.com/spf13/cobra"
)

// serverurlCmd represents the serverurl command
var serverurlCmd = &cobra.Command{
	Use:   "serverurl",
	Short: "Sets the Pleasant Password server URL for pleasant-cli",
	Long: `Sets the Pleasant Password server URL for pleasant-cli
It must be specified as <PROTOCOL>://<URL>(:<PORT>).

If the port is either 80 or 443, it can be inferred from the protocol and can be omitted.

Example:
pleasant-cli config serverurl <SERVER URL>`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		err := pleasant.WriteConfigFile(cfgFile, "ServerUrl", args[0])
		if err != nil {
			pleasant.ExitFatal(err)
		}

		pleasant.Exit("Server URL saved to:", cfgFile)
	},
}

func init() {
	configCmd.AddCommand(serverurlCmd)
}
