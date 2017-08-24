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
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// helmCmd represents the helm command
var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Use Kubernetes Helm with a Kraken cluster",
	Long: `Use Kubernetes Helm with the Kraken
	cluster configured by the specified yaml file`,
	PreRunE: preRunGetKrakenConfig,
	RunE:    run,
}

func run(cmd *cobra.Command, args []string) error {
	cli, backgroundCtx, err := pullKrakenContainerImage(containerImage)

	minorMajorVersion, err := getK8sVersion(cli, backgroundCtx, args)
	if err != nil {
		return err
	}

	helmPath := "/opt/cnct/kubernetes/" + minorMajorVersion + "/bin/helm"

	// Run helm if available, or get user input if no helm available.
	verfiedHelmPath, err := verifyHelmPath(helmPath, cli)
	if err != nil {
		return err
	}

	if strings.Contains(verfiedHelmPath, minorMajorVersion) {
		status, err := runHelm(helmPath, cli, backgroundCtx, args)
		if err != nil {
			return err
		}

		ExitCode = status
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("No version of helm available for Kubernetes " + minorMajorVersion)

		printHelmVersion, err := latestHelmVersion(cli, backgroundCtx, args)
		if err != nil {
			return err
		}

		fmt.Printf("Would you like to run the latest version of helm %s? [Y/N]: ", printHelmVersion)

		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("Fatal: the following error was thrown while reading user input for helm options: %v", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" {
			runHelm(helmPath, cli, backgroundCtx, args)
		} else if response == "n" {
			fmt.Println("No version of Helm running")
		}

		ExitCode = 0
	}

	return nil
}

// Check to see if path exists, else get latest.
func verifyHelmPath(helmPath string, cli *client.Client) (string, error) {
	command := []string{"test", "-f", helmPath}
	ctx, cancel := getTimedContext()
	defer cancel()

	_, statusCode, timeout, err := containerAction(cli, ctx, command, krakenlibConfigPath)
	if err != nil {
		return "", err
	}

	defer timeout()

	// Unless command returns 0 (filepath exists), assign path to latest.
	if statusCode != 0 {
		helmPath = "/opt/cnct/kubernetes/latest/bin/helm"
	}
	return helmPath, nil
}

// Get the k8s version from Krakenlib
func getK8sVersion(cli *client.Client, backgroundCtx context.Context, args []string) (string, error) {
	command := []string{"/kraken/bin/max_k8s_version.sh", krakenlibConfigPath}
	for _, element := range args {
		command = append(command, strings.Split(element, " ")...)
	}

	ctx, cancel := getTimedContext()
	defer cancel()
	resp, _, timeout, err := containerAction(cli, ctx, command, krakenlibConfigPath)
	if err != nil {
		return "", err
	}

	defer timeout()

	out, err := printContainerLogs(cli, resp, backgroundCtx)
	if err != nil {
		return "", err
	}

	s := string(out)
	return s, nil
}

// If no valid helm version, let user know the latest helm version available.
func latestHelmVersion(cli *client.Client, backgroundCtx context.Context, args []string) (string, error) {
	command := []string{"printenv", "K8S_HELM_VERSION_LATEST"}
	for _, element := range args {
		command = append(command, strings.Split(element, " ")...)
	}

	ctx, cancel := getTimedContext()
	defer cancel()
	resp, _, timeout, err := containerAction(cli, ctx, command, krakenlibConfigPath)
	if err != nil {
		return "", err
	}

	defer timeout()

	out, err := printContainerLogs(cli, resp, backgroundCtx)
	if err != nil {
		return "", err
	}

	s := string(out)
	return s, nil
}

// Run helm if valid path or if user wants to run latest helm.
func runHelm(helmPath string, cli *client.Client, backgroundCtx context.Context, args []string) (int, error) {
	path, err := verifyHelmPath(helmPath, cli)
	if err != nil {
		return -1, err
	}

	command := []string{path}
	for _, element := range args {
		command = append(command, strings.Split(element, " ")...)
	}

	ctx, cancel := getTimedContext()
	defer cancel()
	resp, statusCode, timeout, err := containerAction(cli, ctx, command, krakenlibConfigPath)
	if err != nil {
		return -1, err
	}
	defer timeout()

	if out, err := printContainerLogs(cli, resp, backgroundCtx); err != nil {
		return -1, err
	} else {
		fmt.Printf("%s", out)
	}

	return statusCode, nil
}

func init() {
	toolCmd.AddCommand(helmCmd)
}
