package cmd

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func preRunGetKrakenConfig(cmd *cobra.Command, args []string) error {
	if !cmd.Flag("config").Changed {
		fmt.Printf("config file path not given, using default config file location (%s)\n", k2ConfigPath)
	}

	_, err := os.Stat(k2ConfigPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("File %s does not exist!", k2ConfigPath)
	}

	if err != nil {
		return err
	}

	if err := initKrakenClusterConfig(k2ConfigPath); err != nil {
		return err
	}

	return nil
}


func pullKrakenContainerImage(containerImage string) (*client.Client, context.Context, error) {
	terminalSpinner.Prefix = fmt.Sprintf("Pulling image '%s' ",containerImage)
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