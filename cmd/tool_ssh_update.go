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

// sshUpdateCmd represents the update command
var sshUpdateCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh ssh host list",
	Long: `Refresh a list of SSH hosts for an existing Kraken
	cluster configured by the specified yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		spinnerPrefix := fmt.Sprintf("Refreshing ssh config file for cluster '%s' ", getFirstClusterName())

		command := []string{
			"ansible-playbook",
			"-i",
			"ansible/inventory/localhost",
			"ansible/up.yaml",
			"--extra-vars",
			"config_path=" + ClusterConfigPath + " config_base=" + outputLocation + " kraken_action=up ",
			"--tags",
			"ssh_only",
		}

		onFailure := func(out []byte) {
			fmt.Println("ERROR refreshing ssh inventory for " + getFirstClusterName())
			fmt.Printf("%s", out)
			clusterHelpError(Created, ClusterConfigPath)
		}

		onSuccess := func(out []byte) {
			fmt.Println("Done.")
			if logSuccess {
				fmt.Printf("%s", out)
			}
			clusterHelp(Created, ClusterConfigPath)
		}

		ExitCode, err = runKrakenLibCommand(spinnerPrefix, command, ClusterConfigPath, onFailure, onSuccess)
		return err
	},
}

func init() {
	sshCmd.AddCommand(sshUpdateCmd)
}
