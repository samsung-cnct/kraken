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
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// helmCmd represents the helm command
var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Use Kubernetes Helm with K2 cluster",
	Long: `Use Kubernetes Helm with the  K2
	cluster configured by the specified yaml file`,
	PreRunE: preRunE,
	Run:     run,
}

func preRunE(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(k2Config); os.IsNotExist(err) {
		return errors.New("File " + k2Config + " does not exist!")
	}

	initK2Config(k2Config)

	return nil
}

func run(cmd *cobra.Command, args []string) {
	cli := getClient()

	backgroundCtx := getContext()
	pullImage(cli, backgroundCtx, getAuthConfig64(cli, backgroundCtx))
	minorMajorVersion := getK8sVersion(cli, backgroundCtx, args)
	helmPath := "/opt/cnct/kubernetes/" + minorMajorVersion + "/bin/helm"

	// Run helm if available, or get user input if no helm available.
	if strings.Contains(verifyHelmPath(helmPath, cli), minorMajorVersion) {
		runHelm(helmPath, cli, backgroundCtx, args)
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("No version of helm available for Kubernetes " + minorMajorVersion)
		fmt.Printf("Would you like to run the latest version of helm %s? [Y/N]: ", latestHelmVersion(cli, backgroundCtx, args))

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Errorf("Fatal: the following error was thrown while reading user input for helm options: %v", err)
		}
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" {
			runHelm(helmPath, cli, backgroundCtx, args)
		} else if response == "n" {
			fmt.Println("No version of Helm running")
		}
	}
}

// Check to see if path exists, else get latest.
func verifyHelmPath(helmPath string, cli *client.Client) string {
	command := []string{"test", "-f", helmPath}
	ctx, cancel := getTimedContext()
	defer cancel()
	_, statusCode, timeout := containerAction(cli, ctx, command, k2Config)
	defer timeout()

	// Unless command returns 0 (filepath exists), assign path to latest.
	if statusCode != 0 {
		helmPath = "/opt/cnct/kubernetes/latest/bin/helm"
	}
	return helmPath
}

// Get the k8s version from k2
func getK8sVersion(cli *client.Client, backgroundCtx context.Context, args []string) string {
	command := []string{"/kraken/bin/max_k8s_version.sh", k2Config}
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

// If no valid helm version, let user know the latest helm version available.
func latestHelmVersion(cli *client.Client, backgroundCtx context.Context, args []string) string {
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

// Run helm if valid path or if user wants to run latest helm.
func runHelm(helmPath string, cli *client.Client, backgroundCtx context.Context, args []string) int {
	path := verifyHelmPath(helmPath, cli)
	command := []string{path}
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

func init() {
	toolCmd.AddCommand(helmCmd)
}
