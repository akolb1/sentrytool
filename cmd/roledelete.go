// Copyright © 2016 Alex Kolbasov
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var roleDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm"},
	Short:   "delete Sentry roles",
	Long:    `delete Sentry roles.`,
	Run:     roleDelete,
}

func roleDelete(cmd *cobra.Command, args []string) {
	client, err := getClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	verbose := viper.GetBool(verboseOpt)

	roles, _, err := getRoles(cmd, args, true, client)
	if err != nil {
		fmt.Println(toAPIError(err))
		return
	}

	for _, roleName := range roles {
		force, _ := cmd.Flags().GetBool(forceOpt)
		if !force && !askYN(fmt.Sprintf("delete '%s'? ", roleName)) {
			continue
		}
		err = client.RemoveRole(roleName)
		if err != nil {
			fmt.Println(toAPIError(err))
			continue
		}
		if verbose {
			fmt.Println("removed role ", roleName)
		}
	}
}

// askYN prompts a user for a yes/no answer and returns true if user replies
// with anything starting with 'y'
func askYN(prompt string) bool {
	var response string
	fmt.Print(prompt)
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	response = strings.ToLower(response)
	return strings.HasPrefix(response, "y")
}

func init() {
	roleDeleteCmd.Flags().StringP(matchOpt, "m", "", "regexp matching role")
	roleDeleteCmd.Flags().BoolP(forceOpt, "", false, "force deletion")
	roleCmd.AddCommand(roleDeleteCmd)
}
