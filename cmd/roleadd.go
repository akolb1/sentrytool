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
var roleAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add group or permission to a role",
	Long:  `add group or permission to a role.`,
}

func init() {
	roleAddCmd.PersistentFlags().StringP("role", "r", "", "target role")
	roleCmd.AddCommand(roleAddCmd)
}
