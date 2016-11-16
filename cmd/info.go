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
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:           "info",
	Short:         "Print out cluster state information",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long: `Output some basic information on the current 
	cluster state configured by the specified k2 yaml`,
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
		clusterHelpError(Created, k2ConfigPath)
		ExitCode = 0
	},
}

func init() {
	clusterCmd.AddCommand(infoCmd)
}
