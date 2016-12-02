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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// groupCmd represents the role command
var groupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"g"},
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set(verboseOpt, true)
		return listGroups(cmd, args)
	},
	Short: "list, add or remove groups",
	Long: `group command manages Sentry groups. A group can be added to a role or removed from a role.
A single group can belong to multiple roles.

The role is specified with either '-r' flag ro as the first parameter.
The remaining parameters are group names.

Without subcommands lists groups.`,
	Example: `
  sentrytool group list
  sentrytool group grant -r admin_role admin_group finance_group
  sentrytool group grant admin_role finance_group`,
}

func init() {
	// ALl privilege commands operate on a role which can be supplied with -r flag
	groupCmd.PersistentFlags().StringP("role", "r", "", "role name")
	RootCmd.AddCommand(groupCmd)
}
