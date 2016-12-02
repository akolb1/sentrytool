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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//  groupAddCmd represents the role command
var groupAddCmd = &cobra.Command{
	Use:     "grant",
	Aliases: []string{"add"},
	RunE:    addGroupsToRole,
	Short:   "grant group to a role",
	Long: `Grant command associates group with a specific role.
A role should be either specified with -role flag or be the first argument
followed by list of groups.

If -role flag is specified, arguments are group names to add.`,

	Example: `
  # Grant group to a role
  sentrytool group grant -r admin_role admin_group finance_group
  sentrytool group grant admin_role finance_group

  # Revoke group from role
  sentrytool group revoke -r admin_role admin_group`,
}

// addGroupToRole adds a set of groups to the specific role
func addGroupsToRole(cmd *cobra.Command, args []string) error {
	// Get role name
	roleName, _ := cmd.Flags().GetString("role")
	if len(args) == 0 || (roleName == "" && len(args) == 1) {
		return errors.New("missing group name")
	}

	groups := args
	if roleName == "" {
		roleName = args[0]
		groups = args[1:]
	}

	// Get Thrift client
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer client.Close()

	// Verify that roleName is valid
	isValid, err := isValidRole(client, roleName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !isValid {
		return fmt.Errorf("role %s doesn't exist", roleName)
	}

	// Add groups to the role
	if err = client.AddGroupsToRole(roleName, groups); err != nil {
		fmt.Println(err)
		return nil
	}

	verbose := viper.Get(verboseOpt).(bool)
	if verbose {
		listGroups(cmd, groups)
	}

	return nil
}

func init() {
	groupCmd.AddCommand(groupAddCmd)
}
