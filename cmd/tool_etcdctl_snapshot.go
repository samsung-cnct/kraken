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
	path2 "path"

	"github.com/spf13/cobra"
)

var snapshotSaveFile string

var etcdCtlSnapshotCmd = &cobra.Command{
	Use:   "snapshot [save, restore, status]",
	Short: "Run etcdctl snapshot [save, restore, status], primarily designed to be used with a Kraken cluster when using v3 apis",
	Long: `Use etcd's etcdctl snapshot with the Kraken cluster configured by the specified yaml file`,
	PreRunE: preRunEtcdctl,
	RunE: func(cmd *cobra.Command, args []string) error {
		// verify apiversion user selected is supported
		if apiVersion != 3 {
			return fmt.Errorf("The api-version %d does not snapshot v3 apis. Please consider using version 3.", apiVersion)
		}
		_, volumes, flags := etcdctlSnapshotArguments(apiVersion)

		if verbosity {
			fmt.Printf("args passed: %v \n", args)
		}

		_, err := runEtcdSSHDockerCommand("snapshot", flags, []string {args[0], snapshotSaveFile}, volumes)
		return err
	},
	SilenceUsage: true,
	SilenceErrors: true,
}

func init() {
	etcdCtlCmd.AddCommand(etcdCtlSnapshotCmd)

	etcdCtlSnapshotCmd.PersistentFlags().StringVar(
		&snapshotSaveFile,
		"snapshot-file",
		"/home/core/snapshots/snapshot.db",
		"directory currently used to store data.",
	)
}

func etcdctlSnapshotArguments(apiVersion int) ([] string, []string, []string) {
	flags := []string{}
	volumes := []string{}
	environmentVars := []string{}

	if snapshotSaveFile != "" {
		snapshotVolume := path2.Dir(snapshotSaveFile)
		volumes = append(volumes, fmt.Sprintf("%s:%s", snapshotVolume, snapshotVolume))
	}

	environmentVars = append(environmentVars, fmt.Sprintf("%s=%d", envVarETCDCTL_API, apiVersion))
	return environmentVars, volumes, flags
}


