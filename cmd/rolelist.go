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

	"strings"

	"github.com/akolb1/sentrytool/sentryapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var roleListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list roles",
	Long: `list all roles.
If optional '-m regexp' flag is specified, only list roles matching regexp.
If '-g group' flag is specifiedm only list roles for the speecified group.
In verbose mode prints groups for each role.

If roles are provided in arguments, list those roles, otherwise list all matching roles.

Example:

  sentrytool role list -v
  sentrytool role list -v role1 role2
  sentrytool role list -g group1

`,
	Run: listRoles,
}

type roleGroupMap map[string][]string

// getRoles returns a list of Sentry roles from the server
// if useMatcher is true, filters result according to '-m' flag
// If any roles are specified in args, use them
func getRoles(cmd *cobra.Command,
	args []string,
	useMatcher bool,
	client sentryapi.ClientAPI) ([]string, roleGroupMap, error) {
	group, _ := cmd.Flags().GetString(groupOpt)
	roles, roleGroups, err := client.ListRoleByGroup(group)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(roles)

	var matchRegex *regexp.Regexp

	if useMatcher {
		match, _ := cmd.Flags().GetString(matchOpt)
		if match != "" {
			r, err := regexp.Compile(match)
			if err != nil {
				return nil, nil,
					fmt.Errorf("invalid match expression: %s", err)
			}
			matchRegex = r
		}
	}

	result := make([]string, 0, len(roles))
	roleResult := make(roleGroupMap)

	for _, role := range roleGroups {
		// If match is specified, filter by matches
		if matchRegex != nil && !matchRegex.MatchString(role.Name) {
			continue
		}
		result = append(result, role.Name)
		roleResult[role.Name] = role.Groups
	}

	// If there are any roles on the command line, remove any roles that are
	// not mentioned there
	if len(args) == 0 {
		return result, roleResult, nil
	}

	result = make([]string, 0, len(args))
	roleMap := make(roleGroupMap)

	for _, name := range args {
		if r, ok := roleResult[name]; ok {
			result = append(result, name)
			roleMap[name] = r
		}
	}
	return result, roleMap, nil
}

func listRoles(cmd *cobra.Command, args []string) {
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	roles, roleGroups, err := getRoles(cmd, args, true, client)
	if err != nil {
		fmt.Println(err)
		return
	}

	verbose := viper.Get(verboseOpt).(bool)

	for _, r := range roles {
		if !verbose {
			fmt.Println(r)
			continue
		}
		groups := strings.Join(roleGroups[r], ",")
		fmt.Printf("%s: (%s)\n", r, groups)
	}
}

func init() {
	roleListCmd.Flags().StringP(matchOpt, "m", "", "regexp matching role")
	roleListCmd.Flags().StringP(groupOpt, "g", "", "group for a role")
	roleCmd.AddCommand(roleListCmd)
}
