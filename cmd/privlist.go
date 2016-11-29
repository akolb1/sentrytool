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

var privListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"show"},
	RunE:    listPriv,
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

		if len(privList) == 0 {
			fmt.Println(roleName)
			continue
		}

		privs := make([]string, 0, len(privList))
		for _, priv := range privList {
			privs = append(privs, displayPrivilege(roleName, priv))
		}
		fmt.Println(roleName, "=", strings.Join(privs, ", "))

	}
	return nil
}

func displayPrivilege(role string, privilege *sentryapi.Privilege) string {
	parts := []string{}
	if privilege.Server != "" {
		parts = append(parts, "server="+privilege.Server)
	}
	if privilege.Database != "" {
		parts = append(parts, "db="+privilege.Database)
	}
	if privilege.Table != "" {
		parts = append(parts, "table="+privilege.Table)
	}
	if privilege.Column != "" {
		parts = append(parts, "table="+privilege.Column)
	}
	if privilege.URI != "" {
		parts = append(parts, "uri="+privilege.URI)
	}
	if privilege.Action != "" {
		parts = append(parts, "action="+privilege.Action)
	}

	return strings.Join(parts, "->")
}

func init() {
	privListCmd.Flags().StringP(matchOpt, "m", "", "regexp matching role")
	privListCmd.Flags().StringP(groupOpt, "g", "", "group for a role")

	privCmd.AddCommand(privListCmd)
}
