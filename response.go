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
	var quiet, noPager, asJson, asJsonList *bool
	var cols *[]string
	ProjCmd := &cobra.Command{
		Use:     "response",
		Short:   "Response management",
		Args:    cobra.NoArgs,
		Aliases: []string{"res"},
		Version: "v0.0.0",
	}

	ProjLsCmd := &cobra.Command{
		Use:     "list",
		Short:   "List all responses",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Version: "v0.0.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			paginate := !(*noPager)
			buf := bytes.NewBuffer(nil)
			if *asJsonList {
				_ = json.NewEncoder(buf).Encode(inso.Responses)
				return Output(buf, paginate)
			}
			if *asJson {
				for _, rs := range inso.Responses {
					_ = json.NewEncoder(buf).Encode(rs)
				}
				return Output(buf, paginate)
			}
			if *quiet {
				for _, rs := range inso.Responses {
					buf.WriteString(rs.ID + "\n")
				}
				return Output(buf, paginate)
			}
			t := Table{}
			for _, rs := range inso.Responses {
				r := map[string]string{
					"ID":         rs.ID,
					"StatusCode": fmt.Sprint(rs.StatusCode),
					"Modified":   fmt.Sprint(rs.Modified),
					"ParentID":   rs.ParentID,
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
		Short:   "Get response info",
		Long:    "Get basic json info about the given response",
		Version: "v0.0.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			paginate := !(*noPager)
			id := args[0]
			var res *insomnium.Response
			for _, rs := range inso.Responses {
				if rs.ID == id {
					res = &rs
					break
				}
			}
			if res == nil {
				return errors.New("response not found")
			}
			buf := bytes.NewBuffer(nil)
			err := json.NewEncoder(buf).Encode(res)
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
	asJson = ProjLsCmd.PersistentFlags().BoolP("as-json", "j", false, "As json")
	asJsonList = ProjLsCmd.PersistentFlags().BoolP("as-json-list", "J", false, "As json list")
	cols = ProjLsCmd.PersistentFlags().StringSliceP("cols", "c", []string{"ID", "StatusCode", "Modified", "ParentID"}, "Set columns")

	rootCmd.AddCommand(ProjCmd)
}
