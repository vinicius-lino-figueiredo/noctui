package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

func init() {
	var quiet, noPager, asJSON, asJSONList *bool
	var cols *[]string
	WsCmd := &cobra.Command{
		Use:     "workspace",
		Short:   "Workspace management",
		Args:    cobra.NoArgs,
		Aliases: []string{"ws", "wrk"},
		Version: "v0.0.0",
	}

	WsLsCmd := &cobra.Command{
		Use:     "list",
		Short:   "List all workspaces",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Version: "v0.0.0",
		RunE: func(_ *cobra.Command, _ []string) error {
			paginate := !(*noPager)
			buf := bytes.NewBuffer(nil)
			if *asJSONList {
				_ = json.NewEncoder(buf).Encode(inso.Workspaces)
				return Output(buf, paginate)
			}
			if *asJSON {
				for _, ws := range inso.Workspaces {
					_ = json.NewEncoder(buf).Encode(ws)
				}
				return Output(buf, paginate)
			}
			if *quiet {
				for _, ws := range inso.Workspaces {
					buf.WriteString(ws.ID + "\n")
				}
				return Output(buf, paginate)
			}
			t := Table{}
			for _, ws := range inso.Workspaces {
				r := map[string]string{
					"ID":       ws.ID,
					"Name":     ws.Name,
					"Modified": fmt.Sprint(ws.Modified),
					"ParentID": ws.ParentID,
				}
				t = append(t, r)
			}
			err := t.Render(buf, *cols...)
			if err != nil {
				return err
			}
			return Output(buf, paginate)
		},
	}

	WsGetCmd := &cobra.Command{
		Use:     "get",
		Args:    cobra.ExactArgs(1),
		Short:   "Get workspace info",
		Long:    "Get basic json info about the given workspace",
		Version: "v0.0.0",
		RunE: func(_ *cobra.Command, args []string) error {
			paginate := !(*noPager)
			id := args[0]
			var wrk *insomnium.Workspace
			for _, ws := range inso.Workspaces {
				if ws.ID == id {
					wrk = &ws
					break
				}
			}
			if wrk == nil {
				return errors.New("workspace not found")
			}
			buf := bytes.NewBuffer(nil)
			err := json.NewEncoder(buf).Encode(wrk)
			if err != nil {
				return err
			}
			Output(buf, paginate)
			return nil
		},
	}

	WsCmd.AddCommand(WsLsCmd, WsGetCmd)

	quiet = WsCmd.PersistentFlags().BoolP("quiet", "q", false, "Return id only")
	noPager = WsCmd.PersistentFlags().BoolP("no-pager", "P", false, "Disable pager")
	asJSON = WsLsCmd.PersistentFlags().BoolP("as-json", "j", false, "As json")
	asJSONList = WsLsCmd.PersistentFlags().BoolP("as-json-list", "J", false, "As json list")
	cols = WsLsCmd.PersistentFlags().StringSliceP("cols", "c", []string{"ID", "Name", "Modified", "ParentID"}, "Set columns")

	rootCmd.AddCommand(WsCmd)
}
