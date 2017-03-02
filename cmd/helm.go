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

// helmCmd represents the helm command
var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Use Kubernetes Helm with K2 cluster",
	Long: `Use Kubernetes Helm with the  K2 
	cluster configured by the specified yaml file`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(k2Config); os.IsNotExist(err) {
			return errors.New("File " + k2Config + " does not exist!")
		}

		initK2Config(k2Config)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli := getClient()

		backgroundCtx := getContext()
		pullImage(cli, backgroundCtx, getAuthConfig64(cli, backgroundCtx))

		command := []string{"helm"}
		for _, element := range args {
			command = append(command, strings.Split(element, " ")...)
		}

		ctx, cancel := getTimedContext()
		defer cancel()
		resp, statusCode, timeout := containerAction(cli, ctx, command, k2Config)
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

		fmt.Printf("%s", out)

		ExitCode = statusCode
	},
}

func init() {
	toolCmd.AddCommand(helmCmd)
}
