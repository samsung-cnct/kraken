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
	Use:   "update",
	Short: "Refresh ssh host list",
	Long: `Update a list of SSH hosts for the Kraken 
	cluster configured by the specified yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, backgroundCtx, err := pullKrakenContainerImage(containerImage)
		if err != nil {
			return err
		}

		terminalSpinner.Prefix = "Updating ssh inventory for '" + getContainerName() + "' "
		terminalSpinner.Start()

		command := []string{
			"ansible-playbook",
			"-i",
			"ansible/inventory/localhost",
			"ansible/up.yaml",
			"--extra-vars",
			"config_path=" + args[0] + " config_base=" + outputLocation + " kraken_action=up ",
			"--tags",
			"ssh_only",
		}

		ctx, cancel := getTimedContext()
		defer cancel()
		resp, statusCode, timeout, err := containerAction(cli, ctx, command, args[0])
		if err != nil {
			return err
		}
		defer timeout()

		terminalSpinner.Stop()

		out, err := printContainerLogs(
			cli,
			resp,
			backgroundCtx,
		)
		if err != nil {
			fmt.Println("ERROR updating ssh inventory for " + getContainerName())
			return err
		}

		if statusCode != 0 {
			fmt.Println("ERROR updating ssh inventory for " + getContainerName())
			fmt.Printf("%s", out)
			clusterHelpError(Created, args[0])
		} else {
			fmt.Println("Done.")
			if logSuccess {
				fmt.Printf("%s", out)
			}
			clusterHelp(Created, args[0])
		}

		ExitCode = statusCode
		return nil
	},
}

func init() {
	sshCmd.AddCommand(sshUpdateCmd)
}
