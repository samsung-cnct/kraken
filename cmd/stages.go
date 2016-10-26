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

	"github.com/spf13/cobra"
)

// stagesCmd represents the stages command
var stagesCmd = &cobra.Command{
	Use:   "stages",
	Short: "List of k2 stages",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
Following stages are available:

'all'             - do everything
'dryrun'          - only run stages that do local template generation (no infrasructure modification)
'config'          - configuration loader and dependencies
'common'          - template generation for cloud config common to all node types and dependencies
'fabric'          - template generation for network fabric cloud config and dependencies
'etcd'            - template generation for etcd-specific cloud config and dependencies
'master'          - template generation for master-specific cloud config and dependencies
'node'            - template generation for node-specific cloud config and dependencies
'assembler'       - assemble step for generated templates and dependencies
'provider'        - stand up cloud infrastructure and dependencies
'ssh'             - generator for ssh configuration file and dependencies
'readiness'       - wait for cluster to be ready (for the configured definition of 'ready') and dependencies
'services'        - install configured helm charts and dependencies
'config_only'     - configuration loader
'common_only'     - template generation for cloud config common to all node types
'fabric_only'     - template generation for network fabric cloud config
'etcd_only'       - template generation for etcd-specific cloud config
'master_only'     - template generation for master-specific cloud config
'node_only'       - template generation for node-specific cloud config
'assembler_only'  - assemble step for generated templates
'provider_only'   - stand up cloud infrastructure
'ssh_only'        - generator for ssh configuration file
'readiness_only'  - wait for cluster to be ready (for the configured definition of 'ready')
'services_only'   - install configured helm charts
`)
	},
}

func init() {
	topicCmd.AddCommand(stagesCmd)
}
