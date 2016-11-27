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

	"regexp"

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

func getRoles(cmd *cobra.Command) ([]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	roles, err := client.ListRoleByGroup("")
	if err != nil {
		return nil, err
	}
	sort.Strings(roles)

	var match_regex *regexp.Regexp

	match, _ := cmd.Flags().GetString(matchOpt)
	if match != "" {
		r, err := regexp.Compile(match)
		if err != nil {
			return nil, fmt.Errorf("invalid match expression: %s", err)
		}
		match_regex = r
	}

	result := make([]string, 0, len(roles))

	for _, role := range roles {
		// If match is specified, filter by matches
		if match_regex != nil && !match_regex.MatchString(role) {
			continue
		}
		result = append(result, role)
	}
	return result, nil
}

func listRoles(cmd *cobra.Command, args []string) {
	roles, err := getRoles(cmd)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, r := range roles {
		fmt.Println(r)
	}
}

func init() {
	roleCmd.AddCommand(listCmd)
}
