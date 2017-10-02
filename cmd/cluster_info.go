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

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:           "info",
	Short:         "Print out cluster state information",
	SilenceErrors: true,
	SilenceUsage:  false,
	Long: `Output some basic information on the current
	cluster state configured by the specified Krakenlib yaml`,
	PreRunE: preRunGetClusterConfig,
	RunE: func(cmd *cobra.Command, args []string) error {

		// we do not support any additional arguments, we error out then if there are.
		if len(args) > 0 {
			return fmt.Errorf("Unexpected argument(s) passed %v", args)
		}

		clusterHelp(HelpTypeCreated, ClusterConfigPath)
		ExitCode = 0
	},
}

func init() {
	clusterCmd.AddCommand(infoCmd)
}
