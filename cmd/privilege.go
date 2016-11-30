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
)

var privCmd = &cobra.Command{
	Use:     "privilege",
	Aliases: []string{"priv", "p"},
	Short:   "privilege operations",
	RunE:    listPriv,
	Long: `privilege operations: list, grant or revoke privileges.
Argument is a list of roles.

Examples:

  sentrytool privilege create -r r1 -s server1 -d db2 -t table1 -c columnt1 \
      -a insert
  sentrytool privilege list r1
r1 = db=db1->action=all, \
     server=server1->db=db2->table=table1->column=column1->action=insert
`,
}

func init() {
	privCmd.PersistentFlags().StringP("action", "a", "", "action")
	privCmd.PersistentFlags().StringP("server", "s", "", "server name")
	privCmd.PersistentFlags().StringP("database", "d", "", "database ame")
	privCmd.PersistentFlags().StringP("table", "t", "", "table name")
	privCmd.PersistentFlags().StringP("column", "c", "", "column name")
	privCmd.PersistentFlags().StringP("uri", "u", "", "URI")
	privCmd.PersistentFlags().StringP("scope", "", "", "Scope")
	privCmd.PersistentFlags().StringP("service", "", "", "service name")
	privCmd.PersistentFlags().StringP("role", "r", "", "role name")

	privCmd.PersistentFlags().BoolP("grantoption", "", false, "grantOption")

	RootCmd.AddCommand(privCmd)
}
