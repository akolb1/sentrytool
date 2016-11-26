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

// roleCmd represents the role command
var roleCmd = &cobra.Command{
	Use:   "role",
	Aliases: []string{"r", "rl"},
	Short: "Sentry roles manipulation",
	Long: `Create, list or delete roles.
Multiple roles can be created or removed with a single command.`,
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.Get(hostOpt).(string)
		port := viper.Get(portOpt).(int)
		user := viper.Get(userOpt).(string)
		if err := roleList(host, port, user); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(roleCmd)
}