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
	"github.com/spf13/cobra"
)

var apiVersion int
var snapshotSaveFile string
var hostName string
var hostAddr string
var hostPort int
var hostUsername string

var etcdVersion string
var etcdEndpoints string

var etcdCtlCmd = &cobra.Command{
	Use:     "etcdctl",
	Short:   "Run etcdctl, primarily designed to be used with a Kraken cluster",
	Long:    `Use etcd's etcdctl with the Kraken cluster configured by the specified yaml file`,
	PreRunE: preRunGetClusterConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if !isApiVersionSupported(apiVersion) {
			return fmt.Errorf("the api-version %d is not supported", apiVersion)
		}

		// check if save or not
		if ok, file := checkIfSnapshotSubcommand("save", args); ok && file == "" {
			snapshotSaveFile = defaultKrakenEtcdSnapshotFile
			args = append(args, snapshotSaveFile)
		} else if ok {
			snapshotSaveFile = file
		}

		// check if restore or not
		if ok, file := checkIfSnapshotSubcommand("restore", args); ok && file == "" {
			snapshotSaveFile = defaultKrakenEtcdSnapshotFile
			args = append(args, snapshotSaveFile)
		} else if ok {
			snapshotSaveFile = file
		}

		// check if status or not
		// check if save or not
		if ok, file := checkIfSnapshotSubcommand("status", args); ok && file == "" {
			snapshotSaveFile = defaultKrakenEtcdSnapshotFile
			args = append(args, snapshotSaveFile)
		} else if ok {
			snapshotSaveFile = file
		}

		if verbosity {
			fmt.Printf("Arguments passed: %v \n", args)
		}

		ExitCode, err = runEtcdCtl(args)

		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	toolCmd.AddCommand(etcdCtlCmd)

	etcdCtlCmd.PersistentFlags().IntVar(
		&apiVersion,
		"api-version",
		3,
		"Select the api-version to use with etcdctl",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&etcdVersion,
		"etcd-version",
		defaultEtcdDockerVersion,
		"version of etcdctl to run.",
	)

	etcdCtlCmd.PersistentFlags().StringVar(
		&hostName,
		"host-name",
		"",
		"hostname of cluster node to contact, to be used when an ssh-config file is used.",
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
		&etcdEndpoints,
		"endpoints",
		"",
		"list of etcd endpoints, usually list in the form of https://<private-ip>:2379, https://<private-ip>:4001.",
	)
}
