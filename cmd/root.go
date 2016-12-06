// Copyright © 2016 Samsung CNCT
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
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

var cfgFile string
var containerImage string
var outputLocation string
var dockerHost string
var actionTimeout int
var ExitCode int
var keepAlive bool
var logPath string
var logSuccess bool

// progress spinner
var terminalSpinner = spinner.New(spinner.CharSets[35], 200*time.Millisecond)

// init the k2 config viper instance
var clusterConfig = viper.New()

// init the k2cli config viper instance
var k2cliConfig = viper.New()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "k2cli",
	Short: "CLI for k2 Kubernetes cluster provisioner",
	Long: `k2 cli is a command line interface for k2 
	kubernetes cluster provisioner. k2 documentation is available at:
	https://github.com/samsung-cnct/k2`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(outputLocation); os.IsNotExist(err) {
			os.Mkdir(outputLocation, 0755)
		}
	},
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
	cobra.OnInitialize(initK2CliConfig)

	RootCmd.SetHelpCommand(helpCmd)

	// Global flags
	RootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"k2config",
		"k",
		"",
		"config file for k2cli (default \""+os.Getenv("HOME")+"/.k2cli.yaml)\"")
	RootCmd.PersistentFlags().StringVarP(
		&containerImage,
		"image",
		"i",
		"quay.io/samsung_cnct/k2:latest",
		"k2 container image")
	RootCmd.PersistentFlags().StringVarP(
		&outputLocation,
		"output",
		"o",
		os.Getenv("HOME")+"/.kraken",
		"k2 output folder")
	RootCmd.PersistentFlags().StringVarP(
		&dockerHost,
		"docker-host",
		"d",
		"unix:///var/run/docker.sock",
		"docker host address")
	RootCmd.PersistentFlags().IntVarP(
		&actionTimeout,
		"timeout",
		"t",
		1200,
		"timeout (in seconds) for container actions")
	RootCmd.PersistentFlags().BoolVarP(
		&keepAlive,
		"keep-alive",
		"v",
		false,
		"keep stopped containers.")
	RootCmd.PersistentFlags().StringVarP(
		&logPath,
		"log-path",
		"w",
		"",
		"Save output output of container action to path")
	RootCmd.PersistentFlags().BoolVarP(
		&logSuccess,
		"log-success",
		"x",
		false,
		"Display full action logs on success")
}

// initConfig reads in config file and ENV variables if set.
func initK2CliConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		k2cliConfig.SetConfigFile(cfgFile)
	}

	k2cliConfig.SetConfigName(".k2cli") // name of config file (without extension)
	k2cliConfig.AddConfigPath("$HOME")  // adding home directory as first search path
	k2cliConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	k2cliConfig.SetEnvPrefix("k2cli") // prefix for env vars to configure client itself
	k2cliConfig.AutomaticEnv()        // read in environment variables that match

	// If a config file is found, read it in.
	if err := k2cliConfig.ReadInConfig(); err == nil {
		fmt.Println("Using k2cli config file:", k2cliConfig.ConfigFileUsed())
	}
}

func initK2Config(k2config string) {
	clusterConfig.SetConfigFile(k2config)
	clusterConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	clusterConfig.SetEnvPrefix("k2") // prefix for env vars to configure cluster
	clusterConfig.AutomaticEnv()     // read in environment variables that match

	// If a config file is found, read it in.
	if err := clusterConfig.ReadInConfig(); err == nil {
		fmt.Println("Using k2 config file:", clusterConfig.ConfigFileUsed())
	}
}
