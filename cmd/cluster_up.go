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

var upStagesList string

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:           "up [path to K2 config file]",
	Short:         "create a K2 cluster",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long:          `Creates a K2 cluster described in the specified configuration yaml`,
	PreRunE:       preRunGetKrakenConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, backgroundCtx, err := pullKrakenContainerImage(containerImage)
		if err != nil {
			return err
		}

		if verbosity == false {
			terminalSpinner.Prefix = "Bringing up cluster '" + getContainerName() + "' "
			terminalSpinner.Start()
		}

		command := []string{
			"ansible-playbook",
			"-i",
			"ansible/inventory/localhost",
			"ansible/up.yaml",
			"--extra-vars",
			"config_path=" + k2ConfigPath + " config_base=" + outputLocation + " config_forced=" + strconv.FormatBool(configForced) + " kraken_action=up ",
			"--tags",
			upStagesList,
		}

		ctx, cancel := getTimedContext()
		defer cancel()

		resp, statusCode, timeout, err := containerAction(cli, ctx, command, k2ConfigPath)
		if err != nil {
			return err
		}

		defer timeout()

		if verbosity == false {
			terminalSpinner.Stop()
		}

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
			fmt.Println("ERROR bringing up " + getContainerName())
			fmt.Printf("%s", out)
			clusterHelpError(Created, k2ConfigPath)
		} else {
			fmt.Println("Done.")
			if logSuccess {
				fmt.Printf("%s", out)
			}
			clusterHelp(Created, k2ConfigPath)
		}

		ExitCode = statusCode
		return nil
	},
}

func init() {
	clusterCmd.AddCommand(upCmd)
	upCmd.PersistentFlags().StringVarP(
		&upStagesList,
		"stages",
		"s",
		"all",
		"comma-separated list of K2 stages to run. Run 'k2cli help topic stages' for more info.")
}
