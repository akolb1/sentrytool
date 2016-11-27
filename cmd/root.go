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
	defaultThriftPort = 8038
	hostOpt           = "host"
	portOpt           = "port"
	userOpt           = "username"
	matchOpt          = "match"
	forceOpt          = "force"
	componentOpt      = "component"
	verboseOpt        = "verbose"
)

var (
	cfgFile string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sentrytool",
	Short: "Command-line interface to Apache Sentry",
	Long: `Command-line interface to Apache Sentry.
See examples.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	currentUser, _ := user.Current()

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sentrytool.yaml)")
	RootCmd.PersistentFlags().StringP(hostOpt, "H", "localhost", "hostname for Sentry server")
	RootCmd.PersistentFlags().IntP(portOpt, "P", defaultThriftPort, "port for Sentry server")
	RootCmd.PersistentFlags().StringP(userOpt, "U", currentUser.Username, "user name")
	RootCmd.PersistentFlags().StringP(componentOpt, "C", "", "sentry client component")
	RootCmd.PersistentFlags().BoolP(verboseOpt, "v", false, "verbose mode")

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
	viper.SetEnvPrefix("sentry")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//
	}
}
