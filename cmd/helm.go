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
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
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

		//check to see if path exists, else get latest
		helmPath := func() string {
			command := []string{"test", "-f", "/opt/cnct/kubernetes/" + getMinorMajorVersion() + "/bin/helm"}
			ctx, cancel := getTimedContext()
			defer cancel()
			_, statusCode, timeout := containerAction(cli, ctx, command, k2Config)
			defer timeout()

			// unless command return 0 (filepath exists), assign path to latest
			path := "/opt/cnct/kubernetes/latest/bin/helm"
			if statusCode == 0 {
				path = "/opt/cnct/kubernetes/" + getMinorMajorVersion() + "/bin/helm"
			}
			return path
		}

		// Grab latest helm version to let user know which one it is
		//this only happens when there is no valid version of helm for k8s
		latestHelmVersion := func() string {
			command := []string{"printenv", "K8S_HELM_VERSION_LATEST"}
			for _, element := range args {
				command = append(command, strings.Split(element, " ")...)
			}

			ctx, cancel := getTimedContext()
			defer cancel()
			resp, _, timeout := containerAction(cli, ctx, command, k2Config)
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

			s := string(out)
			return s
		}

		//Run helm if valid path or if user wants to run latest helm
		runHelm := func() int {
			command := []string{helmPath()}
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
			return ExitCode

		}

		// Get user input if no helm for k8s version
		if strings.Contains(helmPath(), getMinorMajorVersion()) {
			runHelm()
		} else {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("No version of helm available for Kubernetes " + getMinorMajorVersion())
			fmt.Printf("Would you like to run the latest version of helm %s? [Y/N]: ", latestHelmVersion())

			response, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			response = strings.ToLower(strings.TrimSpace(response))

			if response == "y" {
				runHelm()
			} else if response == "n" {
				fmt.Println("No version of Helm running")
			}
		}
	},
}

func init() {
	toolCmd.AddCommand(helmCmd)
}
