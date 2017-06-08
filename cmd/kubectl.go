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
	"regexp"
	"strings"
)

var kubernetesVersion string

//  This is REALLY bad.  Do NOT do things like this unless compelled (and then be sure to complain)
func fetchNodePools() ([]interface{}, error) {
	clusters_i := clusterConfig.Get("deployment.clusters")
	clusters, okay := clusters_i.([]interface{})
	if !okay {
		fmt.Println(clusters_i)
		return nil, errors.New("Cluster conversion to array failed.")
	}

	c0_i := clusters[0]
	c0, okay := c0_i.(map[interface{}]interface{})
	if !okay {
		fmt.Println(c0_i)
		return nil, errors.New("Cluster 0 conversion to map failed.")
	}

	nodePools_i := c0["nodePools"]
	if nodePools_i == nil {
		fmt.Println(c0)
		return nil, errors.New("nodePool not found for cluster 0")
	}
	nodePools, okay := nodePools_i.([]interface{})
	if nodePools == nil {
		fmt.Println(nodePools_i)
		return nil, errors.New("nodePools conversion failed")
	}
	return nodePools, nil
}

func fetchMasterNodePool() (masterNodePool map[interface{}]interface{}, err error) {
	nodePools, err := fetchNodePools()
	if err != nil {
		return nil, err
	}

	for idx, nodePool_i := range nodePools {
		nodePool, okay := nodePool_i.(map[interface{}]interface{})
		if !okay {
			fmt.Println(nodePool_i)
			return nil, errors.New("nodePool " + string(idx) + " in cluster 0 conversion to map failed.")
		}
		if nodePool["name"] == "master" {
			masterNodePool = nodePool
			break
		}
	}

	if masterNodePool == nil {
		fmt.Println(nodePools)
		return nil, errors.New("Master node pool not found for cluster 0")
	}

	return masterNodePool, nil
}

func setMasterKubectlVersion() error {
	masterNodePool, err := fetchMasterNodePool()
	if err != nil {
		return err
	}

	kubeConfig_i := masterNodePool["kubeConfig"]
	if kubeConfig_i == nil {
		fmt.Println(masterNodePool)
		return errors.New("Master node pool for cluster 0 lacks kubeConfig")
	}
	kubeConfig, okay := kubeConfig_i.(map[interface{}]interface{})
	if !okay {
		fmt.Println(kubeConfig_i)
		return errors.New("Cluster 0 master nodePool kubeConfig conversion to map failed.")
	}

	version_i := kubeConfig["version"]
	if version_i == nil {
		fmt.Println(kubeConfig)
		return errors.New("Cluster 0 master nodePool kubeConfig lacks version")
	}

	version, okay := version_i.(string)
	if !okay {
		fmt.Println(kubeConfig)
		return errors.New("Cluster 0 master nodePool kubeConfig version conversion to string failed.")
	}

	okay, _ = regexp.MatchString(`^v[0-9]+\.[0-9]+\.[0-9]+$`, version)
	if !okay {
		fmt.Println(version)
		return errors.New("Cluster 0 master nodePool kubeConfig version string invalid.")
	}

	re := regexp.MustCompile(`^v[0-9]+\.[0-9]+`)
	// Not checking for okay -- this was implied by validation
	kubernetesVersion = re.FindString(version)

	if kubernetesVersion == "" {
		return errors.New("kubeConfig version empty or not found for first cluster master nodepool.")
	}

	return nil

}

// kubectlCmd represents the kubectl command
var kubectlCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "Use Kubernetes kubectl with K2 cluster",
	Long: `Use Kubernetes kubectl with the K2 
	cluster configured by the specified yaml file`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(k2Config); os.IsNotExist(err) {
			return errors.New("File " + k2Config + " does not exist!")
		}

		initK2Config(k2Config)

		err := setMasterKubectlVersion()
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli := getClient()

		backgroundCtx := getContext()
		pullImage(cli, backgroundCtx, getAuthConfig64(cli, backgroundCtx))

		command := []string{"/opt/cnct/kubernetes/" + kubernetesVersion + "/bin/kubectl"}
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
			getContext(),
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
	toolCmd.AddCommand(kubectlCmd)
}
