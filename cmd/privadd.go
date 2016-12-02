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

var privAddCmd = &cobra.Command{
	Use:     "grant",
	Aliases: []string{"add", "create"},
	Short:   "grant privileges to a role",
	RunE:    addPrivilege,
	Long: `Grant one or several privileges to a role. Privileges can be specified either using
options or using sentry-style privilege specification. Any specification in the command-line
override options.

Multiple privileges may be set at the same time.`,
	Example: `
  $ sentrytool privilege grant -s server2 -r admin \
    'db=db4->table=mytable->action=insert' \
    'db=db5->table=mytable->action=remove'

  $ sentrytool privileges list
  admin = server=server2->db=db4->table=mytable->action=insert,\
          server=server2->db=db5->table=mytable->action=remove`,
}

func addPrivilege(cmd *cobra.Command, args []string) error {
	// Get role name
	roleName, _ := cmd.Flags().GetString("role")
	if roleName == "" && len(args) == 0 {
		return errors.New("missing role name")
	}

	privs := args
	if roleName == "" {
		roleName = args[0]
		privs = args[1:]
	}

	// Get Thrift client
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

	// Add a single privilege from flags

	action, _ := cmd.Flags().GetString("action")
	server, _ := cmd.Flags().GetString("server")
	database, _ := cmd.Flags().GetString("database")
	table, _ := cmd.Flags().GetString("table")
	column, _ := cmd.Flags().GetString("column")
	uri, _ := cmd.Flags().GetString("uri")
	scope, _ := cmd.Flags().GetString("scope")
	grant, _ := cmd.Flags().GetBool("grantoption")
	service, _ := cmd.Flags().GetString("service")
	unsetGrant, _ := cmd.Flags().GetBool("unsetgrant")

	priv := &sentryapi.Privilege{
		Action:           action,
		Server:           server,
		Database:         database,
		Table:            table,
		Column:           column,
		URI:              uri,
		Scope:            scope,
		GrantOption:      grant,
		Service:          service,
		UnsetGrantOption: unsetGrant,
	}

	addPrivileges(client, roleName, priv, privs)
	return nil
}

// Add multiple privileges to a role
func addPrivileges(client sentryapi.ClientAPI, role string, template *sentryapi.Privilege,
	args []string) {
	// Without args, the template is our privilege
	if len(args) == 0 {
		err := client.GrantPrivilege(role, template)
		if err != nil {
			fmt.Println(err)
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
		err = client.GrantPrivilege(role, privilege)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func init() {
	privAddCmd.Flags().BoolP("unsetgrant", "", false, "set grant option to 'unset")
	privCmd.AddCommand(privAddCmd)
}
