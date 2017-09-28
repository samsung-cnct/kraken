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

var updateNodepools string
var addNodepools string
var rmNodepools string

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:           "update [path to kraken config file]",
	Short:         "update a Kraken cluster",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long:          `Updates a Kraken cluster described in the specified configuration yaml`,
	PreRunE: preRunGetClusterConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if updateNodepools == "" && addNodepools == "" && rmNodepools == "" {
			return fmt.Errorf("You must specify which nodepools you want to update. Please pass a comma-separated list of nodepools to update-nodepools, " +
				"add-nodepools or rm-nodepools depending on what action you are taking against the nodepools.  For example: \n kraken cluster update " +
				"--update-nodepools masterNodes,clusterNodes,otherNodes --rm-nodepools badNodepool")
		}

		spinnerPrefix := fmt.Sprintf("Updating cluster '%s' ", getFirstClusterName())

		command := []string{
			"ansible-playbook",
			"-i",
			"ansible/inventory/localhost",
			"ansible/update.yaml",
			"--extra-vars",
			fmt.Sprintf("config_path=%s config_base=%s config_forced=%t kraken_action=update update_nodepools=%s add_nodepools=%s remove_nodepools=%s", ClusterConfigPath, outputLocation, configForced, updateNodepools, addNodepools, rmNodepools),
		}

		onFailure := func(out []byte) {
			fmt.Printf("ERROR updating cluster %s \n", getFirstClusterName())
			fmt.Printf("%s", out)
			clusterHelpError(HelpTypeUpdated, ClusterConfigPath)
		}

		onSuccess := func(out []byte) {
			fmt.Println("Done.")
			if logSuccess {
				fmt.Printf("%s", out)
			}
			clusterHelp(HelpTypeUpdated, ClusterConfigPath)
		}

		ExitCode, err = runKrakenLibCommand(spinnerPrefix, command, ClusterConfigPath, onFailure, onSuccess)
		return err
	},
}

func init() {
	clusterCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().StringVarP(
		&updateNodepools,
		"update-nodepools",
		"",
		"",
		"specify a comma separated list of nodepools to update")
	updateCmd.PersistentFlags().StringVarP(
		&addNodepools,
		"add-nodepools",
		"",
		"",
		"specify a comma separated list of nodepools to add")
	updateCmd.PersistentFlags().StringVarP(
		&rmNodepools,
		"rm-nodepools",
		"",
		"",
		"specify a comma separated list of nodepools to remove")

}
