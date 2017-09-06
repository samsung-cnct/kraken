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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var generatePath string
var provider string
var configPath string

var generateCmd = &cobra.Command{
	Use:          "generate [path to save the Kraken config file at] (default ) " + os.ExpandEnv("$HOME/.kraken/config.yaml"),
	Short:        "Generate a Kraken config file",
	SilenceUsage: true,
	Long:         `Generate a Kraken configuration file at the specified location`,
	PreRunE:      preRunEFunc,
	RunE:         runFunc,
}

func preRunEFunc(cmd *cobra.Command, args []string) error {
	switch provider {
	case "gke":
		configPath = "ansible/roles/kraken.config/files/gke-config.yaml"
	case "aws":
		configPath = "ansible/roles/kraken.config/files/config.yaml"
	default:
		return fmt.Errorf("Error unsupported provider: %s", provider)

	}

	if len(args) > 0 {
		generatePath = os.ExpandEnv(args[0])
	} else {
		generatePath = os.ExpandEnv("$HOME/.kraken/config.yaml")
	}


	fmt.Printf("Attempting to generate configuration at: %s \n", generatePath)

	if _, err := os.Stat(generatePath); !os.IsNotExist(err) {
		return fmt.Errorf("Attempted to create %s, but the file already exists: rename, delete or move it, then run 'generate' subcommand again to generate a new default Kraken config file", generatePath)
	}

	// needed to put define correct location to generate to.
	outputLocation = filepath.Dir(generatePath)

	return os.MkdirAll(outputLocation, 0777)

}

func runFunc(cmd *cobra.Command, args []string) error {
	var err error
	spinnerPrefix := fmt.Sprintf("Generating cluster config %s ", getFirstClusterName())


	command := []string{
		"bash",
		"-c",
		fmt.Sprintf("cp %s %s", configPath, generatePath),
	}


	onFailure := func(out []byte) {
		fmt.Println("Error generating config at " + generatePath)
		fmt.Printf("%s", out)
	}

	onSuccess := func(out []byte) {
		fmt.Printf("Generated %s config at %s \n", provider, generatePath)
		if logSuccess {
			fmt.Printf("%s", out)
		}
	}

	ExitCode, err = runKrakenLibCommand(spinnerPrefix, command, "", onFailure, onSuccess)
	return err
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
