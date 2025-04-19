package main

import (
	"fmt"
	tui "main/internal/tui"

	tv "github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

func init() {
	var workspace *string
	TuiCmd := &cobra.Command{
		Use: "tui",
		RunE: func(_ *cobra.Command, _ []string) error {
			app := tv.NewApplication()
			scr := tui.NewScreen(app, inso)

			if *workspace != "" {
				wrk, found := getWrk(inso, *workspace)
				if !found {
					return fmt.Errorf("workspace %s not found", *workspace)
				}
				scr.OpenWorkspace(wrk)
				return scr.Run("wrk")
			}

			return scr.Run("projs")
		},
	}
	workspace = TuiCmd.PersistentFlags().String("wrk", "", "Workspace")
	rootCmd.AddCommand(TuiCmd)
}

// getWrk receives an id and returns the workspace object and true. If it cannot
// find any workspace with that id, returns nil and false.
func getWrk(inso *insomnium.Insomnium, id string) (*insomnium.Workspace, bool) {
	for _, w := range inso.Workspaces {
		if w.ID == id {
			return &w, true
		}
	}
	return nil, false
}
