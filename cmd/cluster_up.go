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

var upStagesList string

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:           "up [path to Kraken config file]",
	Short:         "Creates a Kraken cluster",
	Long:          `Creates a Kraken cluster described in the specified configuration yaml`,
	PreRunE:       preRunGetClusterConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		spinnerPrefix := fmt.Sprintf("Bringing up cluster '%s' ", getFirstClusterName())

		command := []string{
			"ansible-playbook",
			"-i",
			"ansible/inventory/localhost",
			"ansible/up.yaml",
			"--extra-vars",
			fmt.Sprintf("config_path=%s config_base=%s config_forced=%t kraken_action=up", ClusterConfigPath, outputLocation, configForced),
			"--tags",
			upStagesList,
		}

		onFailure := func(out []byte) {
			fmt.Printf("ERROR bringing up %s \n", getFirstClusterName())
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
	clusterCmd.AddCommand(upCmd)
	upCmd.PersistentFlags().StringVarP(
		&upStagesList,
		"stages",
		"s",
		"all",
		"comma-separated list of Kraken stages to run. Run 'kraken help topic stages' for more info.")
}
