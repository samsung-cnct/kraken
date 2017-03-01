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
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var downStagesList string

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:           "down [path to K2 config file]",
	Short:         "destroy a K2 cluster",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long:          `Destroys a K2 cluster described in the specified configuration yaml`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		k2ConfigPath = os.ExpandEnv("$HOME/.kraken/config.yaml")
		if len(args) > 0 {
			k2ConfigPath = os.ExpandEnv(args[0])
		}

		_, err := os.Stat(k2ConfigPath)
		if os.IsNotExist(err) {
			return errors.New("File " + k2ConfigPath + " does not exist!")
		}

		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		initK2Config(k2ConfigPath)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Pulling image '" + containerImage + "' ")
		terminalSpinner.Start()

		cli := getClient()

		backgroundCtx := getContext()
		pullImage(cli, backgroundCtx, getAuthConfig64(cli, backgroundCtx))

		terminalSpinner.Stop()
		fmt.Println("")

		fmt.Printf("Bringing down cluster '" + getContainerName() + "' ")
		terminalSpinner.Start()

		command := []string{
			"ansible-playbook",
			"-i",
			"ansible/inventory/localhost",
			"ansible/down.yaml",
			"--extra-vars",
			"config_path=" + k2ConfigPath + " config_base=" + outputLocation + " kraken_action=down ",
			"--tags",
			downStagesList,
		}

		ctx, cancel := getTimedContext()
		defer cancel()
		resp, statusCode, timeout := containerAction(cli, ctx, command, k2ConfigPath)
		defer timeout()

		terminalSpinner.Stop()
		fmt.Println("")

		out, err := printContainerLogs(
			cli,
			resp,
			backgroundCtx,
		)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		if len(strings.TrimSpace(logPath)) > 0 {
			writeLog(logPath, out)
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
	},
}

func init() {
	clusterCmd.AddCommand(downCmd)
	downCmd.PersistentFlags().StringVarP(
		&downStagesList,
		"stages",
		"s",
		"all",
		"comma-separated list of K2 stages to run: [all, dryrun, config, fabric, master, node, assembler, readiness, services]")
}
