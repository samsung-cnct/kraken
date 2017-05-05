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
	"os"
	"path/filepath"
)

var generatePath string
var provider string
var configPath string

var generateCmd = &cobra.Command{
	Use:          "generate [path to save the K2 config file at] (default ) " + os.ExpandEnv("$HOME/.kraken/config.yaml"),
	Short:        "Generate a K2 config file",
	SilenceUsage: true,
	Long:         `Generate a K2 configuration file at the specified location`,
  PreRunE: preRunEFunc,
  Run: runFunc,
}
func preRunEFunc(cmd *cobra.Command, args []string)  error {
	if len(args) > 0 {
		generatePath = os.ExpandEnv(args[0])
	} else {
		generatePath = os.ExpandEnv("$HOME/.kraken/config.yaml")
	}

	if provider == "gke" {
		configPath = "ansible/roles/kraken.config/files/gke-config.yaml "
	} else {
		configPath = "ansible/roles/kraken.config/files/config.yaml "
	}

	err := os.MkdirAll(filepath.Dir(generatePath), 0777)
	if err != nil {
		return err
	}
	return nil
}

func runFunc(cmd *cobra.Command, args []string) {
		terminalSpinner.Prefix = "Pulling image '" + containerImage + "' "
		terminalSpinner.Start()

		cli := getClient()

		backgroundCtx := getContext()
		pullImage(cli, backgroundCtx, getAuthConfig64(cli, backgroundCtx))

		terminalSpinner.Stop()

		command := []string{
			"bash",
			"-c",
			"cp " + configPath  + generatePath,
		}

		ctx, cancel := getTimedContext()
		defer cancel()

		outputLocation = filepath.Dir(generatePath)
		resp, statusCode, timeout := containerAction(cli, ctx, command, "")
		defer timeout()

		out, err := printContainerLogs(
			cli,
			resp,
			backgroundCtx,
		)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		if statusCode != 0 {
			fmt.Println("Error generating config at " + generatePath)
			fmt.Printf("%s", out)
		} else {
			fmt.Println("Generated " + provider + " config at " + generatePath)
			if logSuccess {
				fmt.Printf("%s", out)
			}
		}

		ExitCode = statusCode
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringVarP(
		&provider,
		"provider",
		"p",
		"aws",
		"specify a provider for config defaults")
}
