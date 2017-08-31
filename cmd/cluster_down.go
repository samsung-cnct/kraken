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

// deprecated
var downStagesList string
var downtagsList string


// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:           "down [path to Kraken config file]",
	Short:         "destroy a Kraken cluster",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long:          `Destroys a Kraken cluster described in the specified configuration yaml`,
	PreRunE:       preRunGetClusterConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		spinnerPrefix := fmt.Sprintf("Bringing down cluster '%s' ", getFirstClusterName())
		var tagList string

		// remove when deprecation is finalized
		if downStagesList == "all" {
			tagList =  downtagsList
		} else {
			tagList = downStagesList
		}

		command := []string{
			"ansible-playbook",
			"-i",
			"ansible/inventory/localhost",
			"ansible/down.yaml",
			"--extra-vars",
			fmt.Sprintf("config_path=%s config_base=%s config_forced=%t kraken_action=down", ClusterConfigPath, outputLocation, configForced),
			"--tags",
			tagList,
		}

		onFailure := func(out []byte) {
			fmt.Printf("ERROR bringing down %s \n", getFirstClusterName())
			fmt.Printf("%s", out)
		}

		onSuccess := func(out []byte) {
			if logSuccess {
				fmt.Printf("%s", out)
			}
			fmt.Println("Done.")
		}

		ExitCode, err = runKrakenLibCommand(spinnerPrefix, command, ClusterConfigPath, onFailure, onSuccess)
		return err
	},
}

func init() {
	clusterCmd.AddCommand(downCmd)

	downCmd.PersistentFlags().StringVar(
		&downtagsList,
		"tags",
		"all",
		"comma-separated list of Kraken stages to run: [all, dryrun, config, fabric, master, node, assembler, readiness, services]",
	)

	downCmd.PersistentFlags().StringVarP(
		&downStagesList,
		"stages",
		"s",
		"all",
		"[DEPRECATED, please use --tags] comma-separated list of Kraken stages to run: [all, dryrun, config, fabric, master, node, assembler, readiness, services]",
	)

}
