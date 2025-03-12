package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

const insoPathEnv = "INSOMNIUM_PATH"

var rootCmd = &cobra.Command{
	Use:     "noctui",
	Version: "v0.0.0",
}

var inso *insomnium.Insomnium

func main() {
	insPath := rootCmd.PersistentFlags().StringP("insomnium-path", "i", "", "Set insomnium derectory used")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		finalPath := *insPath
		if finalPath == "" {
			finalPath = os.Getenv(insoPathEnv)
			if finalPath == "" {
				return fmt.Errorf("environment variable %s is not set\n", insoPathEnv)
			}
		}
		inso = insomnium.NewInsomnium(finalPath)
		err := inso.Load()
		if err != nil {
			return err
		}
		return nil
	}
	rootCmd.Execute()
}
