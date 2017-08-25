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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var downStagesList string

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:           "down [path to Kraken config file]",
	Short:         "destroy a Kraken cluster",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long:          `Destroys a Kraken cluster described in the specified configuration yaml`,
	PreRunE:       preRunGetClusterConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, backgroundCtx, err := pullKrakenContainerImage(containerImage)
		if err != nil {
			return err
		}

		if verbosity == false {
			terminalSpinner.Prefix = "Bringing down cluster '" + getContainerName() + "' "
			terminalSpinner.Start()
		}

		command := []string{
			"ansible-playbook",
			"-i",
			"ansible/inventory/localhost",
			"ansible/down.yaml",
			"--extra-vars",
			"config_path=" + ClusterConfigPath + " config_base=" + outputLocation + " config_forced=" + strconv.FormatBool(configForced) + " kraken_action=down ",
			"--tags",
			downStagesList,
		}

		ctx, cancel := getTimedContext()
		defer cancel()

		resp, statusCode, timeout, err := containerAction(cli, ctx, command, ClusterConfigPath)
		if err != nil {
			return err
		}

		defer timeout()

		if verbosity == false {
			terminalSpinner.Stop()
		}

		out, err := printContainerLogs(
			cli,
			resp,
			backgroundCtx,
		)
		if err != nil {
			fmt.Printf("ERROR bringing down %s \n", getContainerName())
			return err
		}

		if len(strings.TrimSpace(logPath)) > 0 {
			if err := writeLog(logPath, out); err != nil {
				return err
			}

		}

		if statusCode != 0 {
			fmt.Println("ERROR bringing down " + getContainerName())
			fmt.Printf("%s", out)
		} else {
			if logSuccess {
				fmt.Printf("%s", out)
			}
			fmt.Println("Done.")
		}

		ExitCode = statusCode
		return nil
	},
}

func init() {
	clusterCmd.AddCommand(downCmd)
	downCmd.PersistentFlags().StringVarP(
		&downStagesList,
		"stages",
		"s",
		"all",
		"comma-separated list of Kraken stages to run: [all, dryrun, config, fabric, master, node, assembler, readiness, services]")
}
