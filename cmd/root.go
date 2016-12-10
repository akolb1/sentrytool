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
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultThriftPort = "8038"
	hostOpt           = "host"
	portOpt           = "port"
	userOpt           = "username"
	matchOpt          = "match"
	groupOpt          = "group"
	forceOpt          = "force"
	componentOpt      = "component"
	verboseOpt        = "verbose"
	jstackOpt         = "jstack"
)

var (
	cfgFile string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sentrytool",
	Short: "Command-line interface to Apache Sentry",
	Run:   listAll,
	Long: `Command-line interface to Apache Sentry.
See https://github.com/akolb1/sentrytool/blob/master/doc/sentrytool.md for full documentation

Configuration:

The tool can be configured using either command-line flags, environment variables or
a config file. Config file may be in JSON, TOML, YAML, HCL, and
Java properties config files format. B default the file ~/.sentrytool.yaml is used.

 The following environment variables are used:

* SENTRY_HOST:      Sentry server host name or IP address ('host' in the config file)
* SENTRY_PORT:      Listening port for the Sentry server ('port' in the config file)
* SENTRY_USER:      User name on which behalf the request is made ('user' in the config file)
* SENTRY_COMPONENT: Component name (e.g. 'kafka'). ('component' in the config file)
* SENRY_VERBOSE:    Use verbose mode if set ('verbose' in config file)

When a component is specified the tool uses Generic client model, otherwise it uses the
legacy model.
`,
	Example: `
  # Display everything
  $ sentrytool
  [roles]
  admin
  customer
  [groups]
  g1 = admin
  g2 = admin
  g3 = admin
  user_group = customer
  [privileges]

  # List roles
  $ sentrytool role list
  admin
  customer
  # List roles with groups
  $ sentrytool role list -v
  admin: (g1,g2,g3)
  customer: (user_group)

  # Listing groups
  sentrytool group list
  # Grant and revoke groups to roles
  sentrytool group grant -r admin_role admin_group finance_group
  sentrytool group revoke admin_role finance_group

  # Grant and list privileges
  sentrytool privilege grant -r r1 -s server1 -d db2 -t table1 -c columnt1 \
      -a insert
  sentrytool privilege list r1 r1 = db=db1->action=all, \
     server=server1->db=db2->table=table1->column=column1->action=insert`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// listAllCmd shows all roles, groups and privileges
func listAll(cmd *cobra.Command, args []string) {
	viper.Set(verboseOpt, true)
	fmt.Println("[roles]")
	listRoles(cmd, args)
	fmt.Println("[groups]")
	listGroups(cmd, args)
	fmt.Println("[privileges]")
	listPriv(cmd, args)
}

func init() {
	cobra.OnInitialize(initConfig)

	currentUser, _ := user.Current()

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sentrytool.yaml)")
	RootCmd.PersistentFlags().StringP(hostOpt, "H", "localhost", "hostname for Sentry server")
	RootCmd.PersistentFlags().StringP(portOpt, "P", defaultThriftPort, "port for Sentry server")
	RootCmd.PersistentFlags().StringP(userOpt, "U", currentUser.Username, "user name")
	RootCmd.PersistentFlags().StringP(componentOpt, "C", "", "sentry client component")
	RootCmd.PersistentFlags().BoolP(verboseOpt, "v", false, "verbose mode")
	RootCmd.PersistentFlags().BoolP(jstackOpt, "J", false, "show Java stack on for errors")

	// Bind flags to viper variables
	viper.BindPFlags(RootCmd.PersistentFlags())

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".sentrytool") // name of config file (without extension)
	viper.AddConfigPath("$HOME")       // adding home directory as first search path
	viper.SetEnvPrefix("sentry")       // All environment vars should start with SENTRY_
	viper.AutomaticEnv()               // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//
	}
}
