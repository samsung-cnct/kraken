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
	"strings"

	"github.com/spf13/cobra"
)

// kubectlCmd represents the kubectl command
var kubectlCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "Use Kubernetes kubectl with Kraken cluster",
	Long: `Use Kubernetes kubectl with the Kraken
	cluster configured by the specified yaml file`,
	PreRunE: preRunGetClusterConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var command []string

		// to avoid breaking change, can be removed later when no longer supporting v0.2 kraken-lib
		if comp, err := compareReleases(krakenLibTagToSemver(KrakenlibTag), "0.2.0"); err != nil {
			return err
		} else if comp <= 0 {
			command = []string{"/kraken/bin/computed_kubectl.sh", ClusterConfigPath}
		} else {
			command = []string{"/kraken/bin/computed_kubectl.sh", "--config", ClusterConfigPath}

			if verbosity {
				command = append(command, "--verbose")
			}
		}

		for _, element := range args {
			command = append(command, strings.Split(element, " ")...)
		}

		onFailure := func(out []byte) {
			fmt.Printf("%s \n", out)
		}

		onSuccess := func(out []byte) {
			fmt.Printf("%s \n", out)
		}

		ExitCode, err = runKrakenLibCommandNoSpinner(command, ClusterConfigPath, onFailure, onSuccess)

		return err
	},
}

func init() {
	toolCmd.AddCommand(kubectlCmd)
}
