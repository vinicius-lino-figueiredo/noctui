// Package main is the entry point of the Noctui application, a terminal user
// interface (TUI) designed to interact with Insomnium API project files. It
// provides a navigable and interactive environment for exploring projects,
// workspaces, request groups, and individual HTTP requests.
//
// The main package assembles commands and screens built from components in the
// internal/tui package and defines how to load and operate on the Insomnium
// data structures.
//
// This CLI application is modular and interactive, making it easier to
// visualize and manipulate API structures directly from the terminal.
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
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		finalPath := *insPath
		if finalPath == "" {
			finalPath = os.Getenv(insoPathEnv)
			if finalPath == "" {
				return fmt.Errorf("environment variable %s is not set", insoPathEnv)
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
