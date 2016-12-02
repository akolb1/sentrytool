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
	"strings"

	"github.com/akolb1/sentrytool/sentryapi"
	"github.com/spf13/cobra"
)

/*
 * Sentry config syntax allows things like
 *
 * analyst_role = server=server1->db=analyst1, \
 *   server=server1->db=jranalyst1->table=*->action=select
 *   server=server1->uri=hdfs://ha-nn-uri/landing/analyst1
 *
 * We parse these into a structure representation
 */
const (
	sentrySeparator = "->"
	valSeparator    = "="

	serverKey = "server"
	dbKey     = "db"
	tableKey  = "table"
	columnKey = "column"
	uriKey    = "uri"
	actionKey = "action"
	grantKey  = "grantoption"
)

var privCmd = &cobra.Command{
	Use:     "privilege",
	Aliases: []string{"priv", "p"},
	Short:   "privilege operations",
	RunE:    listPriv,
	Long: `privilege operations: list, grant or revoke privileges.
Argument is a list of roles.`,
	Example: `
  sentrytool privilege grant -r r1 -s server1 -d db2 -t table1 -c columnt1 \
      -a insert
  sentrytool privilege list r1
r1 = db=db1->action=all, \
     server=server1->db=db2->table=table1->column=column1->action=insert`,
}

// Parse privilege in Sentry format into a Privilege object
// E.g. server=server1=>db=mydb
func parsePrivilege(priv string,
	template *sentryapi.Privilege) (*sentryapi.Privilege, error) {
	parts := strings.Split(priv, sentrySeparator)
	privilege := *template
	for _, v := range parts {
		splits := strings.Split(v, valSeparator)
		if len(splits) != 2 {
			return nil, fmt.Errorf("invalid perm format for '%s'", v)
		}
		name := splits[0]
		val := splits[1]
		switch name {
		case serverKey:
			privilege.Server = val
		case dbKey:
			privilege.Database = val
		case tableKey:
			privilege.Table = val
		case columnKey:
			privilege.Column = val
		case uriKey:
			privilege.URI = val
		case actionKey:
			privilege.Action = val
		default:
			return nil, fmt.Errorf("invalid scope name %s", name)
		}

		// Special case for internal purposes - allow setting grantOption to "UNSET"
		if name == grantKey && strings.HasPrefix(strings.ToLower(val), "t") {
			privilege.GrantOption = true
			continue
		}
	}
	return &privilege, nil
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
