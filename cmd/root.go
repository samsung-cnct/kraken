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
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var containerImage string
var outputLocation string
var actionTimeout int
var dockerClient DockerClientConfig
var ExitCode int
var keepAlive bool
var logPath string
var logSuccess bool
var verbosity bool
var KrakenlibTag string // this is set via linker flag

// to be used by subcommands using --config
var ClusterConfigPath string

// progress spinner
var terminalSpinner = spinner.New(spinner.CharSets[35], 200*time.Millisecond)

// init the Krakenlib config viper instance
var clusterConfig = viper.New()

// init the Kraken config viper instance
var krakenConfig = viper.New()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kraken",
	Short: "Kraken is an orchestration and cluster level management system for Kubernetes",
	Long:  "Kraken is an orchestration and cluster level management system for Kubernetes that creates a production scale cluster on a range of platforms using its default settings.",
	PreRun: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(outputLocation); os.IsNotExist(err) {
			os.Mkdir(outputLocation, 0755)
		}
	},
	// to discuss if usage silencing should occur, but errors I think are a must.
	//SilenceUsage: true,
	SilenceErrors: true,
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
	cobra.OnInitialize(initKrakenConfig)
	terminalSpinner.FinalMSG = "Complete\n"

	RootCmd.SetHelpCommand(helpCmd)

	if KrakenlibTag == "" {
		KrakenlibTag = "latest"
	}

	// Populate the global with a "vanilla" Docker configuration
	dockerClient = DockerClientConfig{
		DockerHost:       "",
		DockerAPIVersion: DockerAPIVersion,
		TLSEnabled:       false,
		TLSVerify:        false,
		TLSCACertificate: "",
		TLSCertificate:   "",
		TLSKey:           "",
	}
	// Global flags
	RootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"kraken",
		"k",
		os.ExpandEnv("$HOME/.kraken/.kraken/config.yaml"),
		"Path to kraken config file")
	RootCmd.PersistentFlags().StringVarP(
		&containerImage,
		"image",
		"i",
		"quay.io/samsung_cnct/k2:"+KrakenlibTag,
		"Krakenlib container image")
	RootCmd.PersistentFlags().StringVarP(
		&outputLocation,
		"output",
		"o",
		os.Getenv("HOME")+"/.kraken",
		"Kraken output folder")

	// Specify the docker host string; typically unix:///var/run/docker.sock
	RootCmd.PersistentFlags().StringVarP(
		&dockerClient.DockerHost,
		"docker-host",
		"d",
		dockerClient.GetDefaultHost(),
		"Docker host address")

	// Is TLS supported on the API connection?
	RootCmd.PersistentFlags().BoolVar(
		&dockerClient.TLSEnabled,
		"tls",
		dockerClient.GetDefaultTLSVerify(),
		"Use TLS with the remote API")

	// Should TLS attempt to verify the API connection?
	RootCmd.PersistentFlags().BoolVar(
		&dockerClient.TLSVerify,
		"tlsverify",
		dockerClient.GetDefaultTLSVerify(),
		"Use TLS and verify the remote API")

	// Specify trusted CA for TLS certificates
	RootCmd.PersistentFlags().StringVar(
		&dockerClient.TLSCACertificate,
		"tlscacert",
		dockerClient.GetDefaultTLSCACertificate(),
		"Trust certs signed only by this CA")

	// Specify TLS certificate file
	RootCmd.PersistentFlags().StringVar(
		&dockerClient.TLSCertificate,
		"tlscert",
		dockerClient.GetDefaultTLSCertificate(),
		"Path to the TLS certificate file")

	// Specify TLS Key file
	RootCmd.PersistentFlags().StringVar(
		&dockerClient.TLSKey,
		"tlskey",
		dockerClient.GetDefaultTLSKey(),
		"Path to the TLS key file")

	RootCmd.PersistentFlags().IntVarP(
		&actionTimeout,
		"timeout",
		"t",
		1200,
		"timeout (in seconds) for container actions")
	RootCmd.PersistentFlags().BoolVarP(
		&keepAlive,
		"keep-alive",
		"a",
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
	RootCmd.PersistentFlags().BoolVarP(
		&verbosity,
		"verbosity",
		"v",
		false,
		"Verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initKrakenConfig() {

	// first bind flags
	krakenConfig.BindPFlag("kraken.config", RootCmd.Flags().Lookup("kraken"))
	krakenConfig.BindPFlag("container.image", RootCmd.Flags().Lookup("image"))
	krakenConfig.BindPFlag("output.dir", RootCmd.Flags().Lookup("output"))
	krakenConfig.BindPFlag("docker-host", RootCmd.Flags().Lookup("docker-host"))
	krakenConfig.BindPFlag("tls", RootCmd.Flags().Lookup("tls"))
	krakenConfig.BindPFlag("tlsverify", RootCmd.Flags().Lookup("tlsverify"))
	krakenConfig.BindPFlag("tlscacert", RootCmd.Flags().Lookup("tlscacert"))
	krakenConfig.BindPFlag("tlscert", RootCmd.Flags().Lookup("tlscert"))
	krakenConfig.BindPFlag("tlskey", RootCmd.Flags().Lookup("tlskey"))
	krakenConfig.BindPFlag("timeout", RootCmd.Flags().Lookup("timeout"))
	krakenConfig.BindPFlag("keep-alive", RootCmd.Flags().Lookup("keep-alive"))
	krakenConfig.BindPFlag("log-path", RootCmd.Flags().Lookup("log-path"))
	krakenConfig.BindPFlag("log-success", RootCmd.Flags().Lookup("log-success"))
	krakenConfig.BindPFlag("verbosity", RootCmd.Flags().Lookup("verbosity"))

	// Then env. variables
	krakenConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	krakenConfig.SetEnvPrefix("kraken") // prefix for env vars to configure client itself
	krakenConfig.AutomaticEnv()         // read in environment variables that match

	// Then configs
	krakenConfig.SetConfigName("kraken.config") // name of config file (without extension)
	krakenConfig.AddConfigPath("$HOME/.kraken") // adding home directory as first search path
	krakenConfig.AddConfigPath(".")             // optionally look for config in the working directory

	cfgFile = krakenConfig.GetString("kraken.config")
	if cfgFile != "" { // enable ability to specify config file via flag
		krakenConfig.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := krakenConfig.ReadInConfig(); err == nil {
		fmt.Printf("Using kraken config file: %s \n", krakenConfig.ConfigFileUsed())
	}
}

func initClusterConfig(ClusterConfigPath string) error {
	clusterConfig.SetConfigFile(ClusterConfigPath)
	clusterConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	clusterConfig.SetEnvPrefix("krakenlib") // prefix for env vars to configure cluster
	clusterConfig.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := clusterConfig.ReadInConfig(); err != nil {
		return err
	}

	fmt.Printf("Using Kraken config file: %s \n", clusterConfig.ConfigFileUsed())
	return nil
}
