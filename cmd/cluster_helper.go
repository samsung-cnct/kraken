package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// HelpType is an enum type for cluster post processing message handling
type HelpType int

const (
	// HelpTypeCreated is for the case where cluster up was ran
	HelpTypeCreated HelpType = iota
	// HelpTypeDestroyed is for the case where cluster down was ran
	HelpTypeDestroyed
	// HelpTypeUpdated is for the case where cluster udpated was ran
	HelpTypeUpdated
)

func preRunGetClusterConfig(cmd *cobra.Command, args []string) error {
	if !cmd.Flag("config").Changed {
		fmt.Printf("config file path not given, using default config file location (%s)\n", ClusterConfigPath)
	}

	_, err := os.Stat(ClusterConfigPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", ClusterConfigPath)
	}

	if err != nil {
		return err
	}

	if err := initClusterConfig(ClusterConfigPath); err != nil {
		return err
	}

	return nil
}

func pullKrakenContainerImage(containerImage string) (*client.Client, context.Context, error) {
	terminalSpinner.Prefix = fmt.Sprintf("Pulling image '%s' ", containerImage)
	terminalSpinner.Start()

	cli, err := getClient()
	if err != nil {
		return nil, nil, err
	}

	backgroundCtx := getContext()
	authConfig64, err := getAuthConfig64(backgroundCtx, cli)
	if err != nil {
		return nil, nil, err
	}

	if err = pullImage(backgroundCtx, cli, authConfig64); err != nil {
		return nil, nil, err
	}

	terminalSpinner.Stop()
	return cli, backgroundCtx, nil
}

func runKrakenLibCommand(spinnerPrefix string, command []string, clusterConfigPath string, onError func([]byte), onSuccess func([]byte)) (int, error) {
	cli, backgroundCtx, err := pullKrakenContainerImage(containerImage)
	if err != nil {
		return 1, err
	}

	// verbosity false here means show spinner but no container output
	if !verbosity {
		terminalSpinner.Prefix = spinnerPrefix
		terminalSpinner.Start()
	}

	ctx, cancel := getTimedContext()
	defer cancel()

	resp, statusCode, timeout, err := containerAction(ctx, cli, command, clusterConfigPath)
	if err != nil {
		return 1, err
	}

	defer timeout()

	// verbosity false here means show spinner but no container output
	if !verbosity {
		terminalSpinner.Stop()
	}

	out, err := printContainerLogs(backgroundCtx, cli, resp)
	if err != nil {
		return 1, err
	}

	if len(strings.TrimSpace(logPath)) > 0 {
		if err := writeLog(logPath, out); err != nil {
			return 1, err
		}
	}

	if statusCode != 0 {
		onError(out)
	} else {
		onSuccess(out)
	}

	return statusCode, nil
}

func runKrakenLibCommandNoSpinner(command []string, clusterConfigPath string, onError func([]byte), onSuccess func([]byte)) (int, error) {
	cli, backgroundCtx, err := pullKrakenContainerImage(containerImage)
	if err != nil {
		return 1, err
	}

	ctx, cancel := getTimedContext()
	defer cancel()

	resp, statusCode, timeout, err := containerAction(ctx, cli, command, clusterConfigPath)
	if err != nil {
		return 1, err
	}

	defer timeout()

	out, err := printContainerLogs(backgroundCtx, cli, resp)
	if err != nil {
		return 1, err
	}

	if len(strings.TrimSpace(logPath)) > 0 {
		if err := writeLog(logPath, out); err != nil {
			return 1, err
		}
	}

	if statusCode != 0 {
		onError(out)
	} else {
		onSuccess(out)
	}

	return statusCode, nil
}

func clusterHelpError(help HelpType, clusterConfigFile string) {
	switch help {
	case HelpTypeCreated:
		fmt.Printf("ERROR: bringing up cluster %s, using config file %s \n", getFirstClusterName(), clusterConfigFile)
		clusterHelp(help, clusterConfigFile)
	case HelpTypeDestroyed:
		fmt.Printf("ERROR bringing down cluster %s, using config file %s \n", getFirstClusterName(), clusterConfigFile)
		clusterHelp(help, clusterConfigFile)
	case HelpTypeUpdated:
		fmt.Printf("ERROR updating cluster %s, using config file %s \n", getFirstClusterName(), clusterConfigFile)
		clusterHelp(help, clusterConfigFile)
	}

}

func clusterHelp(help HelpType, clusterConfigFile string) {
	// this doesnt have to be a switch statement, but we may handle these errors different later on, so should be.
	switch help {
	case HelpTypeCreated, HelpTypeUpdated, HelpTypeDestroyed:
		fmt.Println("Some of the cluster state MAY be available:")
		if _, err := os.Stat(path.Join(outputLocation, getFirstClusterName(), "admin.kubeconfig")); err == nil {
			fmt.Println("To use kubectl: ")
			fmt.Println(" kubectl --kubeconfig=" + path.Join(
				outputLocation,
				getFirstClusterName(), "admin.kubeconfig") + " [kubectl commands]")
			fmt.Println(" or use 'kraken tool kubectl --config " + clusterConfigFile + " [kubectl commands]'")

			if _, err := os.Stat(path.Join(outputLocation,
				getFirstClusterName(), "admin.kubeconfig")); err == nil {
				fmt.Println("To use helm: ")
				fmt.Println(" export KUBECONFIG=" + path.Join(
					outputLocation,
					getFirstClusterName(), "admin.kubeconfig"))
				fmt.Println(" helm [helm command] --home " + path.Join(
					outputLocation,
					getFirstClusterName(), ".helm"))
				fmt.Println(" or use 'kraken tool helm --config " + clusterConfigFile + " [helm commands]'")
			}
		}

		if _, err := os.Stat(path.Join(outputLocation, getFirstClusterName(), "ssh_config")); err == nil {
			fmt.Println("To use ssh: ")
			fmt.Println(" ssh <node pool name>-<number> -F " + path.Join(outputLocation, getFirstClusterName(), "ssh_config"))
			// This is usage has not been implemented. See issue #49
			//fmt.Println(" or use 'kraken tool --config ssh ssh " + clusterConfigFile + " [ssh commands]'")
		}
	}

}

func getFirstClusterName() string {
	defaultVal := "cluster-name-missing"

	if clusters := clusterConfig.Get("deployment.clusters"); clusters != nil {
		// return the first cluster item in the slice, then cast into an object
		firstCluster := clusters.([]interface{})[0].(map[interface{}]interface{})
		return getConfigValueOrDefault(firstCluster, "name", defaultVal)
	}

	return defaultVal
}

func getFirstClusterKubernetesVersion(provider string) (string, error) {
	defaultValue := "Unknown version"
	defaultError := "error: could not obtain kubernetes version, reason: %s"

	switch provider {
	case "aws":
		clusters := clusterConfig.Get("deployment.clusters")
		if clusters == nil {
			return defaultValue, fmt.Errorf(defaultError, "could not find any cluster entry")
		}

		clustersList := clusters.([]interface{})
		if len(clustersList) == 0 {
			return defaultValue, fmt.Errorf(defaultError, "cluster list is empty.")
		}

		firstCluster := clusters.([]interface{})[0].(map[interface{}]interface{})

		nodePools := firstCluster["nodePools"]
		if nodePools == nil {
			return defaultValue, fmt.Errorf(defaultError, "could not find nodePools entry")
		}

		nodePoolsList := nodePools.([]interface{})
		if len(nodePoolsList) == 0 {
			return defaultValue, fmt.Errorf(defaultError, "nodePool list is empty.")
		}

		var apiserverNode map[interface{}]interface{}

		for _, nodeRaw := range nodePoolsList {
			node := nodeRaw.(map[interface{}]interface{})

			if node != nil && node["apiServerConfig"] != nil {
				apiserverNode = node
				break
			}
		}

		if apiserverNode == nil  {
			return defaultValue, fmt.Errorf(defaultError, "could not apiserver nodes")
		}

		kubeconfig := apiserverNode["kubeConfig"]
		if kubeconfig == nil  {
			return defaultValue, fmt.Errorf(defaultError, "could not find kubeconfig key")
		}

		kubeconfigObj := kubeconfig.(map[interface{}]interface{})
		version := kubeconfigObj["version"]

		if version == nil {
			return defaultValue, fmt.Errorf(defaultError, "could not find kubeconfig version key")
		}

		if verbosity {
			fmt.Printf("found kubernetes version from config file, version is: %s", version.(string))
		}

		return version.(string), nil
	case "gke":
		clusters := clusterConfig.Get("deployment.clusters")
		if clusters == nil {
			return defaultValue, fmt.Errorf(defaultError, "could not find any cluster entry")
		}

		clustersList := clusters.([]interface{})
		if len(clustersList) == 0 {
			return defaultValue, fmt.Errorf(defaultError, "cluster list is empty.")
		}

		firstCluster := clusters.([]interface{})[0].(map[interface{}]interface{})

		nodePools := firstCluster["nodePools"]
		if nodePools == nil {
			return defaultValue, fmt.Errorf(defaultError, "could not find nodePools entry")
		}

		nodePoolsList := nodePools.([]interface{})
		if len(nodePoolsList) == 0 {
			return defaultValue, fmt.Errorf(defaultError, "nodePool list is empty.")
		}

		node := nodePoolsList[0].(map[interface{}]interface{})
		if node == nil  {
			return defaultValue, fmt.Errorf(defaultError, "could not retrive a valid node")
		}

		kubeconfig := node["kubeConfig"]
		if kubeconfig == nil  {
			return defaultValue, fmt.Errorf(defaultError, "could not find kubeconfig key")
		}

		kubeconfigObj := kubeconfig.(map[interface{}]interface{})
		version := kubeconfigObj["version"]

		if version == nil {
			return defaultValue, fmt.Errorf(defaultError, "could not find kubeconfig version key")
		}

		if verbosity {
			fmt.Printf("found kubernetes version from config file, version is: %s", version.(string))
		}

		return version.(string), nil
	}

	return defaultValue, fmt.Errorf(defaultError, "could not find kubernetes version")

}

func getProviderFromConfig() (string, error) {
	defaultValue := "Unknown Provider"
	defaultError := "error: could not obtain kubernetes version, reason: %s"

	clusters := clusterConfig.Get("deployment.clusters")
	if clusters == nil {
		return defaultValue, fmt.Errorf(defaultError, "could not find any cluster entry")
	}

	clustersList := clusters.([]interface{})
	if len(clustersList) == 0 {
		return defaultValue, fmt.Errorf(defaultError, "cluster list is empty.")
	}

	firstCluster := clusters.([]interface{})[0].(map[interface{}]interface{})

	providerConfig := firstCluster["providerConfig"]
	if providerConfig == nil {
		return defaultValue, fmt.Errorf(defaultError, "could not find a providerConfig key")
	}

	providerConfigObj := providerConfig.(map[interface{}]interface{})
	provider := providerConfigObj["provider"]

	if provider == nil {
		return defaultValue, fmt.Errorf(defaultError, "could not find the provider key")
	}

	if verbosity {
		fmt.Printf("found cluster provider from config file, provider is: %s", provider.(string))
	}

	return provider.(string), nil
}


func getConfigValueOrDefault(configObj map[interface{}]interface{}, key string, defaultVal string) string {
	if configObj[key] == nil {
		return defaultVal
	}
	// should not use type assertion .(string) without verifying interface isnt nil
	return os.ExpandEnv(configObj[key].(string))
}


