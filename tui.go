package main

import (
	tui "main/internal/tui"

	tv "github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func init() {

	TuiCmd := &cobra.Command{
		Use: "tui",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := tv.NewApplication()
			scr := tui.NewScreen(app, inso)
			return scr.Run("projs")
		},
	}
	rootCmd.AddCommand(TuiCmd)
}
