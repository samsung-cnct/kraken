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
	"os"

	"github.com/spf13/cobra"
)

// toolCmd represents the tool command
var toolCmd = &cobra.Command{
	Use:           "tool",
	Short:         "Use tools with Kraken cluster",
	SilenceUsage:  true,
	SilenceErrors: true,
	Long: `Use various third-party tools with a
	Kraken cluster configured by specified yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		ExitCode = 0
	},
}

func init() {
	RootCmd.AddCommand(toolCmd)

	toolCmd.PersistentFlags().StringVarP(
		&ClusterConfigPath,
		"config",
		"c",
		os.ExpandEnv("$HOME/.kraken/config.yaml"),
		"required path to the kraken cluster config")
}
