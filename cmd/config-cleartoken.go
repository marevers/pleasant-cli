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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// clearTokenCmd represents the cleartoken command
var clearTokenCmd = &cobra.Command{
	Use:   "cleartoken",
	Short: "Removes the token file",
	Long: `Removes the token file
	
Example:
pleasant-cli config cleartoken`,
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Remove(tokenFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Token file deleted:", tokenFile)
	},
}

func init() {
	configCmd.AddCommand(clearTokenCmd)
}
