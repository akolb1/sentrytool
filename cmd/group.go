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

import "github.com/spf13/cobra"

// roleCmd represents the role command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "list, add or remove groups to the role",
	Long: `List, add or remove groups to the role.
The role is specified with either '-r' flag ro as the first parameter.
The remaining parameters are group names.

Examples:

    group list
    group grant -r admin_role admin_group finance_group
    group grant admin_role finance_group
`,
	RunE:  listGroups,
}

func init() {
	groupCmd.Flags().StringP("role", "r", "", "roleName")
	RootCmd.AddCommand(groupCmd)
}
