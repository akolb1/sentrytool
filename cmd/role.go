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
	"fmt"
)

// roleCmd represents the role command
var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Sentry roles manipulation",
	Long: `Create, list or delete roles.
Multiple roles can be created or removed with a single command.`,
	Example: `
	# List roles
	$ sentrytool role list
	admin
	customer
	# List roles with groups
	$ sentrytool role list -v
	admin: (g1,g2,g3)
	customer: (user_group)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set(verboseOpt, true)
		if len(args) != 0 {
			return fmt.Errorf("invalid subcommand %s", args[1])
		}
		listRoles(cmd, args)
		return nil
	},
}

func init() {
	roleCmd.Flags().StringP(matchOpt, "m", "", "regexp matching role")
	RootCmd.AddCommand(roleCmd)
}
