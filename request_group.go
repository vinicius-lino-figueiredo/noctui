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
	ReqGrpCmd := &cobra.Command{
		Use:     "request-group",
		Short:   "RequestGroup management",
		Args:    cobra.NoArgs,
		Aliases: []string{"rg", "rqgr"},
		Version: "v0.0.0",
	}

	ReqGrpLsCmd := &cobra.Command{
		Use:     "list",
		Short:   "List all request groups",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Version: "v0.0.0",
		RunE: func(_ *cobra.Command, _ []string) error {
			paginate := !(*noPager)
			buf := bytes.NewBuffer(nil)
			if *asJSONList {
				_ = json.NewEncoder(buf).Encode(inso.RequestGroups)
				return Output(buf, paginate)
			}
			if *asJSON {
				for _, rg := range inso.RequestGroups {
					_ = json.NewEncoder(buf).Encode(rg)
				}
				return Output(buf, paginate)
			}
			if *quiet {
				for _, rg := range inso.RequestGroups {
					buf.WriteString(rg.ID + "\n")
				}
				return Output(buf, paginate)
			}
			t := Table{}
			for _, rg := range inso.RequestGroups {
				r := map[string]string{
					"ID":       rg.ID,
					"Name":     rg.Name,
					"Modified": fmt.Sprint(rg.Modified),
					"ParentID": rg.ParentID,
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

	ReqGrpGetCmd := &cobra.Command{
		Use:     "get",
		Args:    cobra.ExactArgs(1),
		Short:   "Get request group info",
		Long:    "Get basic json info about the given request group",
		Version: "v0.0.0",
		RunE: func(_ *cobra.Command, args []string) error {
			paginate := !(*noPager)
			id := args[0]
			var rqgr *insomnium.RequestGroup
			for _, rg := range inso.RequestGroups {
				if rg.ID == id {
					rqgr = &rg
					break
				}
			}
			if rqgr == nil {
				return errors.New("request group not found")
			}
			buf := bytes.NewBuffer(nil)
			err := json.NewEncoder(buf).Encode(rqgr)
			if err != nil {
				return err
			}
			Output(buf, paginate)
			return nil
		},
	}

	ReqGrpCmd.AddCommand(ReqGrpLsCmd, ReqGrpGetCmd)

	quiet = ReqGrpCmd.PersistentFlags().BoolP("quiet", "q", false, "Return id only")
	noPager = ReqGrpCmd.PersistentFlags().BoolP("no-pager", "P", false, "Disable pager")
	asJSON = ReqGrpLsCmd.PersistentFlags().BoolP("as-json", "j", false, "As json")
	asJSONList = ReqGrpLsCmd.PersistentFlags().BoolP("as-json-list", "J", false, "As json list")
	cols = ReqGrpLsCmd.PersistentFlags().StringSliceP("cols", "c", []string{"ID", "Name", "Modified", "ParentID"}, "Set columns")

	rootCmd.AddCommand(ReqGrpCmd)
}
