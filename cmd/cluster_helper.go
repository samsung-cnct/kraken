package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

)

func preRunGetClusterConfig(cmd *cobra.Command, args []string) error {
	if !cmd.Flag("config").Changed {
		fmt.Printf("config file path not given, using default config file location (%s)\n", ClusterConfigPath)
	}

	_, err := os.Stat(ClusterConfigPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("File %s does not exist!", ClusterConfigPath)
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
	authConfig64, err := getAuthConfig64(cli, backgroundCtx)
	if err != nil {
		return nil, nil, err
	}

	if err = pullImage(cli, backgroundCtx, authConfig64); err != nil {
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

	resp, statusCode, timeout, err := containerAction(cli, ctx, command, clusterConfigPath)
	if err != nil {
		return 1, err
	}

	defer timeout()

	// verbosity false here means show spinner but no container output
	if !verbosity {
		terminalSpinner.Stop()
	}

	out, err := printContainerLogs(cli, resp, backgroundCtx)
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

	resp, statusCode, timeout, err := containerAction(cli, ctx, command, clusterConfigPath)
	if err != nil {
		return 1, err
	}

	defer timeout()

	out, err := printContainerLogs(cli, resp, backgroundCtx)
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
