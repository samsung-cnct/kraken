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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		ClusterConfigPath = os.ExpandEnv("$HOME/.kraken/config.yaml")
		if len(args) > 0 {
			ClusterConfigPath = os.ExpandEnv(args[0])
		}

		if updateNodepools == "" && addNodepools == "" && rmNodepools == "" {
			return fmt.Errorf("You must specify which nodepools you want to update. Please pass a comma-separated list of nodepools to update-nodepools, " +
				"add-nodepools or rm-nodepools depending on what action you are taking against the nodepools.  For example: \n kraken cluster update " +
				"--update-nodepools masterNodes,clusterNodes,otherNodes --rm-nodepools badNodepool")
		}

		_, err := os.Stat(ClusterConfigPath)
		if os.IsNotExist(err) {
			return fmt.Errorf("File %s does not exist!", ClusterConfigPath)
		}

		if err != nil {
			return err
		}

		if err := initClusterConfig(ClusterConfigPath); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, backgroundCtx, err := pullKrakenContainerImage(containerImage)
		if err != nil {
			return err
		}

		terminalSpinner.Prefix = "Updating cluster '" + getContainerName() + "' "
		terminalSpinner.Start()

		command := []string{
			"ansible-playbook",
			"-i",
			"ansible/inventory/localhost",
			"ansible/update.yaml",
			"--extra-vars",
			"config_path=" + ClusterConfigPath + " config_base=" + outputLocation + " kraken_action=update " + " update_nodepools=" + updateNodepools +
				" add_nodepools=" + addNodepools + " remove_nodepools=" + rmNodepools,
		}

		ctx := getContext()
		// defer cancel()
		resp, statusCode, timeout, err := containerAction(cli, ctx, command, ClusterConfigPath)
		if err != nil {
			return err
		}

		defer timeout()

		terminalSpinner.Stop()

		out, err := printContainerLogs(cli, resp, backgroundCtx)
		if err != nil {
			return err
		}

		if len(strings.TrimSpace(logPath)) > 0 {
			if err := writeLog(logPath, out); err != nil {
				return err
			}
		}

		if statusCode != 0 {
			fmt.Println("ERROR updating " + getContainerName())
			fmt.Printf("%s", out)
			clusterHelpError(Created, ClusterConfigPath)
		} else {
			fmt.Println("Done.")
			if logSuccess {
				fmt.Printf("%s", out)
			}
			clusterHelp(Created, ClusterConfigPath)
		}

		ExitCode = statusCode
		return nil
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
