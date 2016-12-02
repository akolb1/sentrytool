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
	"github.com/spf13/viper"
)

var privListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"show", "ls"},
	Short:   "list matching privileges",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set(verboseOpt, true)
		return listPriv(cmd, args)
	},
	Long: `list all matching privileges for given roles.
Roles are given as command-line arguments.

If any of the filtering options (server, database, table, etc) are specified,
  only show matching privileges.`,
}

func listPriv(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer client.Close()

	roles, _, err := getRoles(cmd, args, true, client)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// For list command these flags mean filtering options
	action, _ := cmd.Flags().GetString("action")
	server, _ := cmd.Flags().GetString("server")
	database, _ := cmd.Flags().GetString("database")
	table, _ := cmd.Flags().GetString("table")
	column, _ := cmd.Flags().GetString("column")
	uri, _ := cmd.Flags().GetString("uri")
	scope, _ := cmd.Flags().GetString("scope")
	grant, _ := cmd.Flags().GetBool("grantoption")
	service, _ := cmd.Flags().GetString("service")

	for _, roleName := range roles {
		isValid, err := isValidRole(client, roleName)
		if err != nil {
			return err
		}
		if !isValid {
			return fmt.Errorf("role %s doesn't exist", roleName)
		}
		privList, err := client.ListPrivilegesByRole(roleName)
		if err != nil {
			fmt.Println(err)
			continue
		}

		privs := make([]string, 0, len(privList))
		// Go through privileges and add matching ones
		for _, priv := range privList {
			if action != "" && priv.Action != action {
				continue
			}
			if server != "" && priv.Server != server {
				continue
			}
			if database != "" && priv.Database != database {
				continue
			}
			if table != "" && priv.Table != table {
				continue
			}
			if column != "" && priv.Column != column {
				continue
			}
			if uri != "" && priv.URI != uri {
				continue
			}
			if scope != "" && priv.Scope != scope {
				continue
			}
			if service != "" && priv.Service != service {
				continue
			}
			if grant && !priv.GrantOption {
				continue
			}
			privs = append(privs, displayPrivilege(roleName, priv))
		}
		if len(privs) == 0 {
			fmt.Println(roleName)
		} else {
			fmt.Println(roleName, "=", strings.Join(privs, ", "))
		}

	}
	return nil
}

func displayPrivilege(role string, privilege *sentryapi.Privilege) string {
	parts := []string{}
	if privilege.Server != "" {
		parts = append(parts, serverKey+valSeparator+privilege.Server)
	}
	if privilege.Database != "" {
		parts = append(parts, dbKey+valSeparator+privilege.Database)
	}
	if privilege.Table != "" {
		parts = append(parts, tableKey+valSeparator+privilege.Table)
	}
	if privilege.Column != "" {
		parts = append(parts, columnKey+valSeparator+privilege.Column)
	}
	if privilege.URI != "" {
		parts = append(parts, uriKey+valSeparator+privilege.URI)
	}
	if privilege.Action != "" {
		parts = append(parts, actionKey+valSeparator+privilege.Action)
	}

	return strings.Join(parts, sentrySeparator)
}

func init() {
	privListCmd.Flags().StringP(matchOpt, "m", "", "regexp matching role")
	privListCmd.Flags().StringP(groupOpt, "g", "", "group for a role")

	privCmd.AddCommand(privListCmd)
}
