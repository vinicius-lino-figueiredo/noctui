package main

import (
	tv "github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func init() {

	TuiCmd := &cobra.Command{
		Use: "tui",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := tv.NewApplication()
			scr := NewScreen(app)
			return scr.Run("projs")
		},
	}
	rootCmd.AddCommand(TuiCmd)
}
