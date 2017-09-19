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

var etcdDataDir string
var etcdBackupDir string


var etcdCtlBackupCmd = &cobra.Command{
	Use:   "backup [command options] [arguments...]",
	Short: "Run etcdctl backup, primarily designed to be used with a Kraken cluster but should be only used if using v2 apis",
	Long:  "Use etcd's etcdctl backup with the Kraken cluster configured by the specified yaml file",
	PreRunE: preRunEtcdctl,
	RunE: func(cmd *cobra.Command, args []string) error {
		// verify apiVersion user selected is supported
		if apiVersion != 2 {
			return fmt.Errorf("The api-version %d is not meant for v2 apis. Please use version 2.", apiVersion)
		}

		_, volumes, flags := etcdctlBackupArguments(apiVersion)

		_, err := runEtcdSSHDockerCommand("backup", flags, args, volumes)
		return err
	},
	SilenceUsage: true,
	SilenceErrors: true,
}

func init() {
	etcdCtlCmd.AddCommand(etcdCtlBackupCmd)

	etcdCtlBackupCmd.PersistentFlags().StringVar(
		&etcdDataDir,
		"data-dir",
		"/ephemeral/etcd",
		"directory currently used to store data.",
	)

	etcdCtlBackupCmd.PersistentFlags().StringVar(
		&etcdBackupDir,
		"backup-dir",
		"/ephemeral/etcd_backup",
		"directory currently used to store backup data, will be created if does not exist.",
	)

}



func etcdctlBackupArguments(apiVersion int) ([] string, []string, []string) {
	flags := []string{}
	volumes := []string{}
	environmentVars := []string{}

	if etcdDataDir != "" {
		flags = append(flags, fmt.Sprintf("--data-dir %s", etcdDataDir))
		volumes = append(volumes, fmt.Sprintf("%s:%s", etcdDataDir, etcdDataDir))
	}

	if etcdBackupDir != "" {
		flags = append(flags, fmt.Sprintf("--backup-dir %s", etcdBackupDir))
		volumes = append(volumes, fmt.Sprintf("%s:%s", etcdBackupDir, etcdBackupDir))
	}

	environmentVars = append(environmentVars, fmt.Sprintf("%s=%d", envVarETCDCTL_API, apiVersion))

	return environmentVars, volumes, flags
}


