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
)

var userName string
var password string
var logPath string
var logSuccess bool
var k2ConfigPath string

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "k2 cluster actions",
	Long:  `Commands that perform actions on a k2 cluster described by provided yaml config`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		ExitCode = 1
	},
}

func init() {
	RootCmd.AddCommand(clusterCmd)
	clusterCmd.PersistentFlags().StringVarP(
		&userName,
		"user",
		"u",
		"",
		"registry user name")
	clusterCmd.PersistentFlags().StringVarP(
		&password,
		"password",
		"p",
		"",
		"registry password")
	clusterCmd.PersistentFlags().StringVarP(
		&logPath,
		"log-path",
		"w",
		"",
		"Save output output of container action to path")
	clusterCmd.PersistentFlags().BoolVarP(
		&logSuccess,
		"log-success",
		"x",
		false,
		"Display full action logs on success")

}
