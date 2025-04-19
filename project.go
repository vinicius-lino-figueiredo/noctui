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
	ProjCmd := &cobra.Command{
		Use:     "project",
		Short:   "Project management",
		Args:    cobra.NoArgs,
		Aliases: []string{"prj"},
		Version: "v0.0.0",
	}

	ProjLsCmd := &cobra.Command{
		Use:     "list",
		Short:   "List all projects",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Version: "v0.0.0",
		RunE: func(_ *cobra.Command, _ []string) error {
			paginate := !(*noPager)
			buf := bytes.NewBuffer(nil)
			if *asJSONList {
				_ = json.NewEncoder(buf).Encode(inso.Projects)
				return Output(buf, paginate)
			}
			if *asJSON {
				for _, pr := range inso.Projects {
					_ = json.NewEncoder(buf).Encode(pr)
				}
				return Output(buf, paginate)
			}
			if *quiet {
				for _, pr := range inso.Projects {
					buf.WriteString(pr.ID + "\n")
				}
				return Output(buf, paginate)
			}
			t := Table{}
			for _, pr := range inso.Projects {
				r := map[string]string{
					"ID":       pr.ID,
					"Name":     pr.Name,
					"Modified": fmt.Sprint(pr.Modified),
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

	ProjGetCmd := &cobra.Command{
		Use:     "get",
		Args:    cobra.ExactArgs(1),
		Short:   "Get project info",
		Long:    "Get basic json info about the given project",
		Version: "v0.0.0",
		RunE: func(_ *cobra.Command, args []string) error {
			paginate := !(*noPager)
			id := args[0]
			var prj *insomnium.Project
			for _, pr := range inso.Projects {
				if pr.ID == id {
					prj = &pr
					break
				}
			}
			if prj == nil {
				return errors.New("project not found")
			}
			buf := bytes.NewBuffer(nil)
			err := json.NewEncoder(buf).Encode(prj)
			if err != nil {
				return err
			}
			Output(buf, paginate)
			return nil
		},
	}

	ProjCmd.AddCommand(ProjLsCmd, ProjGetCmd)

	quiet = ProjCmd.PersistentFlags().BoolP("quiet", "q", false, "Return id only")
	noPager = ProjCmd.PersistentFlags().BoolP("no-pager", "P", false, "Disable pager")
	asJSON = ProjLsCmd.PersistentFlags().BoolP("as-json", "j", false, "As json")
	asJSONList = ProjLsCmd.PersistentFlags().BoolP("as-json-list", "J", false, "As json list")
	cols = ProjLsCmd.PersistentFlags().StringSliceP("cols", "c", []string{"ID", "Name", "Modified", "ParentID"}, "Set columns")

	rootCmd.AddCommand(ProjCmd)
}
