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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var roleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create Sentry roles",
	Long:  `create Sentry roles.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer client.Close()

		verbose := viper.Get(verboseOpt).(bool)
		for _, roleName := range args {
			err = client.CreateRole(roleName)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if verbose {
				fmt.Println("created role ", roleName)
			}
		}
	},
}

func init() {
	roleCmd.AddCommand(roleCreateCmd)
}
