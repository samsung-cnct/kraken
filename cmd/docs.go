// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var docPath string

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:          "docs [output dir]",
	Short:        "Generate markdown docs",
	SilenceUsage: true,
	Long:         `Generate complete markdown doc tree for k2cli`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			docPath = os.ExpandEnv(args[0])
			err := os.MkdirAll(docPath, 0777)
			if err != nil {
				return err
			}
		} else {
			docPath = "./"
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		doc.GenMarkdownTree(RootCmd, docPath)
	},
}

func init() {
	RootCmd.AddCommand(docsCmd)
}
