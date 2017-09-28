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
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const (
	tmpFile      string = "k8s_version.txt"
	versionRegex string = `(v\d+\.\d+\.\d+)`
)

// helmCmd represents the helm command
var helmCmd = &cobra.Command{
	Use:           "helm",
	Short:         "Use Kubernetes Helm with a Kraken cluster",
	SilenceUsage:  true,
	SilenceErrors: true,
	Long:          `Use Kubernetes Helm with the Kraken cluster configured by the specified yaml file`,
	PreRunE:       preRunGetClusterConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		cli, backgroundCtx, err := pullKrakenContainerImage(containerImage)
		if err != nil {
			return err
		}

		minorMajorVersion, err := getK8sVersion(cli)
		if err != nil {
			return err
		}

		helmPath := path.Join("/opt/cnct/kubernetes/", minorMajorVersion, "/bin/helm")

		// Run helm if available, or get user input if no helm available.
		verifiedHelmPath, err := verifyHelmPath(helmPath, cli)
		if err != nil {
			return err
		}

		if strings.Contains(verifiedHelmPath, minorMajorVersion) {
			ExitCode, err = runHelm(backgroundCtx, helmPath, cli, args)
			return err
		}

		fmt.Printf("No version of helm available for Kubernetes %s \n", minorMajorVersion)

		latestHelmVersion, err := latestSupportedHelmVersion(backgroundCtx, cli)
		if err != nil {
			ExitCode = -1
			return err
		}

		fmt.Printf("Would you like to run the latest version of helm %s [Y/N]?: \n", latestHelmVersion)

		response, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			ExitCode = 1
			return fmt.Errorf("Fatal: the following error was thrown while reading user input for helm options: %v", err)
		}

		switch strings.ToLower(strings.TrimSpace(response)) {
		case "y", "yes":
			ExitCode, err = runHelm(backgroundCtx, helmPath, cli, args)
			return err
		case "n", "no":
			fmt.Println("No version of Helm running")
		default:
			return fmt.Errorf("Please answer with only 'Y' or 'N'")
		}

		ExitCode = 0
		return err
	},
}

func init() {
	toolCmd.AddCommand(helmCmd)
}

// Check to see if path exists, else get latest.
func verifyHelmPath(helmPath string, cli *client.Client) (string, error) {
	command := []string{"test", "-f", helmPath}

	statusCode, err := runContainerCommand(nil, cli, command, nil)

	// Unless command returns 0 (filepath exists), assign path to latest.
	if statusCode != 0 {
		helmPath = "/opt/cnct/kubernetes/latest/bin/helm"
	}

	return helmPath, err
}

// Get the k8s version from Krakenlib
func getK8sVersion(cli *client.Client) (string, error) {
	k8sVersionErr := fmt.Errorf("Error: retrieving k8s version from config file: %s", ClusterConfigPath)

	outputFile := fmt.Sprintf("%s_%s", ClusterConfigPath, tmpFile)
	command := []string{
		"ansible-playbook",
		"-i",
		"ansible/inventory/localhost",
		"ansible/max_k8s_version.yaml",
		"--extra-vars",
		fmt.Sprintf("config_path=%s config_base=%s config_forced=%t kraken_action=max_k8s_version version_outfile=%s", ClusterConfigPath, outputLocation, configForced, outputFile),
	}

	statusCode, err := runContainerCommand(nil, cli, command, nil)
	if err != nil {
		return "", err
	}

	if statusCode != 0 {
		return "", k8sVersionErr
	}

	// now we find the file from the outputFile
	file, err := os.Open(outputFile)
	if err != nil {
		return "", k8sVersionErr
	}

	defer Close(file)
	defer remove(outputFile)

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", k8sVersionErr
	}
	version := strings.TrimSuffix(scanner.Text(), "\n")

	return removePatchVersion(version)
}

// removePatchVersion removes patch number from version string,
func removePatchVersion(version string) (string, error) {
	if verbosity {
		fmt.Printf("The cluster's Kubernetes version is %s.\n", version)
	}

	r, err := regexp.Compile(versionRegex)
	if err != nil {
		return "", err
	}

	err = fmt.Errorf("error: unexpected version: %s", version)

	if !r.MatchString(version) {
		return "", err
	}

	versionArray := strings.Split(version, ".")

	if len(versionArray) != 3 {
		return "", err
	}

	return strings.Join([]string{versionArray[0], versionArray[1]}, "."), nil
}

// Run helm if valid path or if user wants to run latest helm.
func runHelm(backgroundCtx context.Context, helmPath string, cli *client.Client, args []string) (int, error) {
	path, err := verifyHelmPath(helmPath, cli)
	if err != nil {
		return -1, err
	}

	command := []string{path}
	for _, element := range args {
		command = append(command, strings.Split(element, " ")...)
	}

	onComplete := func(out []byte) {
		fmt.Printf("%s", out)
	}

	return runContainerCommand(backgroundCtx, cli, command, onComplete)
}

func remove(path string) {
	if err := os.Remove(path); err != nil {
		log.Fatal(err)
	}
}

// If no valid helm version, let user know the latest helm version available.
func latestSupportedHelmVersion(backgroundCtx context.Context, cli *client.Client) (string, error) {
	var result string

	command := []string{"printenv", "K8S_HELM_VERSION_LATEST"}

	onComplete := func(out []byte) {
		result = string(out)
	}

	_, err := runContainerCommand(backgroundCtx, cli, command, onComplete)

	return result, err
}

func runContainerCommand(backgroundCtx context.Context, cli *client.Client, command []string, onComplete func([]byte)) (int, error) {
	var err error
	ctx, cancel := getTimedContext()

	defer cancel()

	resp, statusCode, timeout, err := containerAction(ctx, cli, command, ClusterConfigPath)
	if err != nil {
		return -1, err
	}

	defer timeout()

	if backgroundCtx != nil {
		out, err := printContainerLogs(backgroundCtx, cli, resp)

		if err != nil {
			return -1, err
		}

		if onComplete != nil {
			onComplete(out)
		}

	}

	return statusCode, err

}
