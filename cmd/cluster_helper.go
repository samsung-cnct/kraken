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
	if ClusterConfigPath == "" {
		return fmt.Errorf("please pass a valid kraken config file")
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
	clusterName := getFirstClusterName()

	switch help {
	case HelpTypeCreated, HelpTypeUpdated, HelpTypeDestroyed:
		fmt.Println("\nSome of the cluster state MAY be available:")

		// output that depends on admin.kubeconfig existing
		kubeConfigPath := path.Join(outputLocation, clusterName, "admin.kubeconfig")
		if _, err := os.Stat(kubeConfigPath); err == nil {
			// kubectl
			fmt.Println("\nTo use kubectl: ")
			fmt.Printf(" kubectl --kubeconfig=%s [kubectl commands]\n", kubeConfigPath)

			if outputLocation == os.ExpandEnv("$HOME/.kraken") {
				fmt.Printf(" or use 'kraken tool --config %s kubectl [kubectl commands]'\n", clusterConfigFile)
			} else {
				fmt.Printf(" or use 'kraken tool --config %s --output %s kubectl [kubectl commands]'\n", clusterConfigFile, outputLocation)
			}

			// helm
			helmPath := path.Join(outputLocation, clusterName, ".helm")
			if _, err := os.Stat(helmPath); err == nil {
				fmt.Println("\nTo use helm: ")
				fmt.Printf(" export KUBECONFIG=%s\n", kubeConfigPath)
				fmt.Printf(" helm [helm command] --home %s\n", helmPath)

				if outputLocation == os.ExpandEnv("$HOME/.kraken") {
					fmt.Printf(" or use 'kraken tool --config %s helm [helm commands]'\n", clusterConfigFile)
				} else {
					fmt.Printf(" or use 'kraken tool --config %s --output %s helm [helm commands]'\n", clusterConfigFile, outputLocation)
				}
			}
		}

		// output that depends on ssh_config existing
		sshConfigPath := path.Join(outputLocation, clusterName, "ssh_config")
		if _, err := os.Stat(sshConfigPath); err == nil {
			// ssh tool
			fmt.Println("\nTo use ssh: ")
			fmt.Printf(" ssh <node pool name>-<number> -F %s\n", sshConfigPath)
			// This is usage has not been implemented. See issue #49
			//fmt.Println(" or use 'kraken tool --config ssh ssh " + clusterConfigFile + " [ssh commands]'")
		}
	}

}

func getFirstClusterName() string {
	// only supports first cluster name right now

	if clusters := clusterConfig.Get("deployment.clusters"); clusters != nil {
		firstCluster := clusters.([]interface{})[0].(map[interface{}]interface{})
		if firstCluster["name"] == nil {
			return "cluster-name-missing"
		}
		// should not use type assertion .(string) without verifying interface isnt nil
		return os.ExpandEnv(firstCluster["name"].(string))
	}

	return "cluster-name-missing"
}
