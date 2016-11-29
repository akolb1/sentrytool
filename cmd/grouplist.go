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

	"github.com/spf13/cobra"
	"sort"
	"strings"
)

var groupListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list groups",
	RunE:    listGroups,
}

// groupMap is a map from group name to list of roles
type groupMap map[string][]string

func listGroups(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer client.Close()

	var matchingGroups map[string]bool
	if len(args) != 0 {
		// Mark all groups mentioned on the command line
		matchingGroups = make(map[string]bool, len(args))
		for _, g := range args {
			matchingGroups[g] = true
		}
	}
	// Get list of all groups and their roles
	_, roleGroups, err := getRoles(cmd, nil, true, client)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	groupMap := make(groupMap)
	for roleName, groups := range roleGroups {
		for _, group := range groups {
			if matchingGroups != nil && !matchingGroups[group] {
				continue
			}
			groupMap[group] = append(groupMap[group], roleName)
		}
	}

	// Get sorted list of groups
	groups := make([]string, 0, len(groupMap))
	for group := range groupMap {
		groups = append(groups, group)
	}
	sort.Strings(groups)

	// Display all groups
	for _, group := range groups {
		fmt.Println(group, "=", strings.Join(groupMap[group], ", "))
	}

	return nil
}

func init() {
	groupCmd.AddCommand(groupListCmd)
}
