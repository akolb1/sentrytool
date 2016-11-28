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
	"github.com/spf13/cobra/doc"
)

// docCmd represents command for writing documentation
var docCmd = &cobra.Command{
	Use:     "doc",
	Aliases: []string{"man"},
	Short:   "write documentation",
	Long:    `Write sentrytool documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		docdir, _ := cmd.Flags().GetString("dir")
		isMan, _ := cmd.Flags().GetBool("man")
		if isMan {
			header := &doc.GenManHeader{
				Title: "Sentry",
				Section: "3",
			}
			doc.GenManTree(RootCmd, header, docdir)
		} else {
			doc.GenMarkdownTree(RootCmd, docdir + "/")
		}
	},
}

func init() {
	docCmd.Flags().StringP("dir", "d", ".", "document directory")
	docCmd.Flags().BoolP("man", "m", false, "generate man pages")
	RootCmd.AddCommand(docCmd)
}
