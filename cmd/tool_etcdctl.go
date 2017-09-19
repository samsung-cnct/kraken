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
	"os"

	"github.com/spf13/cobra"
	"fmt"
)

var apiVersion int
var sshPath string
var sshKey string
var hostAddr string
var hostPort int
var hostUsername string

var etcdVersion string
var etcdCaFile string
var etcdCertFile string
var etcdKeyFile string
var etcdEndpoints string
var etcdUseKrakenCerts bool
var etcdClusterName string



var etcdCtlCmd = &cobra.Command{
	Use:   "etcdctl",
	Short: "Run etcdctl, primarily designed to be used with a Kraken cluster",
	Long: `Use etcd's etcdctl with the Kraken cluster configured by the specified yaml file`,
	PreRunE: preRunEtcdctl,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if verbosity {
			fmt.Printf("Arguments passed: %v \n", args)
		}

		ExitCode, err = runEtcdSSHDockerCommand("", nil, args, nil)

		return err
	},
	SilenceUsage: true,
	SilenceErrors: true,
}

func init() {
	toolCmd.AddCommand(etcdCtlCmd)

	etcdCtlCmd.PersistentFlags().IntVar(
		&apiVersion,
		"api-version",
		2,
		"Select the api-version to use with etcdctl",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&etcdVersion,
		"etcd-version",
		defaultEtcdDockerVersion,
		"version of etcdctl to run.",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&sshKey,
		"ssh-key",
		os.ExpandEnv("$HOME/.ssh/id_rsa"),
		"ssh keyfile to use to ssh into your cluster.",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&hostAddr,
		"host-address",
		"",
		"hostname of cluster node to contact, used in 'ssh username@host:port'.",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&hostUsername,
		"username",
		"core",
		"The username in 'ssh username@host:port' ",
	)

	etcdCtlCmd.PersistentFlags().IntVar(
		&hostPort,
		"port",
		22,
		"port to use in 'ssh username@host:port' ",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&etcdCaFile,
		"ca-file",
		"",
		"location of ca-cert file.",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&etcdCertFile,
		"cert-file",
		"",
		"location of cert file",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&etcdKeyFile,
		"key-file",
		"",
		"location of key file.",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&etcdEndpoints,
		"endpoints",
		"",
		"list of etcd endpoints, usually list in the form of https://<private-ip>:2379, https://<private-ip>:4001.",
	)

	etcdCtlCmd.PersistentFlags().BoolVar(
		&etcdUseKrakenCerts,
		"use-kraken-certs",
		false,
		"set to true if using default kraken ssl certs, do not define 'cluster-name' flag for the single cluster per kraken node case.",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&etcdClusterName,
		"cluster-name",
		"",
		"use only if multiple etcd clusters exist in one kraken node, this defines the cluster's name as defined in kranken config file to use, affects the location used when the flag `use-kraken-certs` is used.",
	)

}

