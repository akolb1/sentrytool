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
	"regexp"
	"sort"

	"github.com/akolb1/sentrytool/sentryapi"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var roleListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list roles",
	Long: `list all roles.
If optional '-m regexp' flag is specified, only list roles matching regexp.`,
	Run: listRoles,
}

// getRoles returns a list of Sentry roles from the server
// if useMatcher is true, filters result according to '-m' flag
func getRoles(cmd *cobra.Command,
	useMatcher bool,
	client sentryapi.SentryClientAPI) ([]string, error) {
	group, _ := cmd.Flags().GetString(groupOpt)
	roles, err := client.ListRoleByGroup(group)
	if err != nil {
		return nil, err
	}
	sort.Strings(roles)

	var matchRegex *regexp.Regexp

	if useMatcher {
		match, _ := cmd.Flags().GetString(matchOpt)
		if match != "" {
			r, err := regexp.Compile(match)
			if err != nil {
				return nil,
					fmt.Errorf("invalid match expression: %s", err)
			}
			matchRegex = r
		}
	}

	result := make([]string, 0, len(roles))

	for _, role := range roles {
		// If match is specified, filter by matches
		if matchRegex != nil && !matchRegex.MatchString(role) {
			continue
		}
		result = append(result, role)
	}
	return result, nil
}

func listRoles(cmd *cobra.Command, _ []string) {
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	roles, err := getRoles(cmd, true, client)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, r := range roles {
		fmt.Println(r)
	}
}

func init() {
	roleListCmd.Flags().StringP(matchOpt, "m", "", "regexp matching role")
	roleListCmd.Flags().StringP(groupOpt, "g", "", "group for a role")
	roleCmd.AddCommand(roleListCmd)
}
