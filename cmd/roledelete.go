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
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		verbose := viper.Get(verboseOpt).(bool)

		var existingRoles map[string]bool
		roles := args
		if len(roles) == 0 {
			r, err := getRoles(cmd, true, client)
			if err != nil {
				log.Fatal(err)
			}
			roles = r
		} else {
			existing, err := getRoles(cmd, false, client)
			if err != nil {
				fmt.Println(err)
				return
			}
			existingRoles = make(map[string]bool)
			for _, role := range existing {
				existingRoles[role] = true
			}
		}

		for _, roleName := range roles {
			if existingRoles != nil && !existingRoles[roleName] {
				fmt.Println("role", roleName, "does not exist, skipping")
				continue
			}

			force, _ := cmd.Flags().GetBool(forceOpt)
			if !force && !askYN(fmt.Sprintf("delete '%s'? ", roleName)) {
				continue
			}
			err = client.RemoveRole(roleName)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if existingRoles != nil {
				existingRoles[roleName] = false
			}
			if verbose {
				fmt.Println("removed role ", roleName)
			}
		}
	},
}

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
