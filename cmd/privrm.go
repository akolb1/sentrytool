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

	"github.com/akolb1/sentrytool/sentryapi"
	"github.com/spf13/cobra"
)

var privRevokeCmd = &cobra.Command{
	Use:     "revoke",
	Aliases: []string{"remove", "delete", "rm"},
	Short:   "revoke privilege",
	Long: `Revoke one or several privileges from a role. Privileges can be specified either using
options or using sentry-style privilege specification. Any specification in the command-line
override options.

Multiple privileges may be set at the same time.`,
	Example: `
  $ sentrytool privilege revoke -s server2 -r admin \
    'db=db4->table=mytable->action=insert' \
    'db=db5->table=mytable->action=remove'`,
	RunE: revokePrivilege,
}

func revokePrivilege(cmd *cobra.Command, args []string) error {
	roleName, _ := cmd.Flags().GetString("role")
	if roleName == "" && len(args) == 0 {
		return errors.New("missing role name")
	}

	privs := args
	if roleName == "" {
		roleName = args[0]
		privs = args[1:]
	}

	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer client.Close()

	isValid, err := isValidRole(client, roleName)
	if err != nil {
		return err
	}
	if !isValid {
		return fmt.Errorf("role %s doesn't exist", roleName)
	}

	action, _ := cmd.Flags().GetString("action")
	server, _ := cmd.Flags().GetString("server")
	database, _ := cmd.Flags().GetString("database")
	table, _ := cmd.Flags().GetString("table")
	column, _ := cmd.Flags().GetString("column")
	uri, _ := cmd.Flags().GetString("uri")
	scope, _ := cmd.Flags().GetString("scope")
	grant, _ := cmd.Flags().GetBool("grantoption")
	service, _ := cmd.Flags().GetString("service")

	priv := &sentryapi.Privilege{
		Action:      action,
		Server:      server,
		Database:    database,
		Table:       table,
		Column:      column,
		URI:         uri,
		Scope:       scope,
		GrantOption: grant,
		Service:     service,
	}

	removePrivileges(client, roleName, priv, privs)
	return nil
}

// Add multiple privileges to a role
func removePrivileges(client sentryapi.ClientAPI, role string, template *sentryapi.Privilege,
	args []string) {
	// Without args, the template is our privilege
	if len(args) == 0 {
		err := client.RevokePrivilege(role, template)
		if err != nil {
			fmt.Println(toAPIError(err))
		}
	}
	// Privileges specified at the command line, parse them, fill unset parts from
	// the template and grant them
	for _, privSpec := range args {
		privilege, err := parsePrivilege(privSpec, template)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = client.RevokePrivilege(role, privilege)
		if err != nil {
			fmt.Println(toAPIError(err))
			continue
		}
	}
}

func init() {
	privCmd.AddCommand(privRevokeCmd)
}
