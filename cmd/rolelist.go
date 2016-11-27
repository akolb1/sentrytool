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

	"sort"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list roles",
	Long:    `list all roles.`,
	Run:     listRoles,
}

func listRoles(cmd *cobra.Command, args []string) {
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	roles, err := client.ListRoleByGroup("")
	if err != nil {
		fmt.Println(err)
		return
	}
	sort.Strings(roles)
	for _, r := range roles {
		fmt.Println(r)
	}
}

func init() {
	roleCmd.AddCommand(listCmd)
}
