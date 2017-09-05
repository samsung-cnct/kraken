package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
)

const(
	krakenImage string = "kraken_store"
	basePath string = "data"
)
var releaseCmd = &cobra.Command{
	Use:   "release-it",
	Short: "Release the kraken!",
	Long: `Tool to release the kraken!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resource, err := Asset(path.Join(basePath, krakenImage))
		if err != nil {
			return err
		}

		str := string(resource) // convert content to a 'string'

		fmt.Println(str)

		ExitCode = 0
		return nil
	},
	SilenceUsage: true,
}


func init() {
	RootCmd.AddCommand(releaseCmd)
}