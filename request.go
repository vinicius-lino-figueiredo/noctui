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
	ReqCmd := &cobra.Command{
		Use:     "request",
		Short:   "Request management",
		Args:    cobra.NoArgs,
		Aliases: []string{"req"},
		Version: "v0.0.0",
	}

	ReqLsCmd := &cobra.Command{
		Use:     "list",
		Short:   "List all requests",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Version: "v0.0.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			paginate := !(*noPager)
			buf := bytes.NewBuffer(nil)
			if *asJsonList {
				_ = json.NewEncoder(buf).Encode(inso.Requests)
				return Output(buf, paginate)
			}
			if *asJson {
				for _, rq := range inso.Requests {
					_ = json.NewEncoder(buf).Encode(rq)
				}
				return Output(buf, paginate)
			}
			if *quiet {
				for _, rq := range inso.Requests {
					buf.WriteString(rq.ID + "\n")
				}
				return Output(buf, paginate)
			}
			t := Table{}
			for _, rq := range inso.Requests {
				r := map[string]string{
					"ID":        rq.ID,
					"Method":    rq.Method,
					"Name":      rq.Name,
					"Modified":  fmt.Sprint(rq.Modified),
					"ParentID":  rq.ParentID,
					"IsPrivate": fmt.Sprint(rq.IsPrivate),
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

	ReqGetCmd := &cobra.Command{
		Use:     "get",
		Args:    cobra.ExactArgs(1),
		Short:   "Get request info",
		Long:    "Get basic json info about the given request",
		Version: "v0.0.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			paginate := !(*noPager)
			id := args[0]
			var req *insomnium.Request
			for _, rq := range inso.Requests {
				if rq.ID == id {
					req = &rq
					break
				}
			}
			if req == nil {
				return errors.New("request not found")
			}
			buf := bytes.NewBuffer(nil)
			err := json.NewEncoder(buf).Encode(req)
			if err != nil {
				return err
			}
			Output(buf, paginate)
			return nil
		},
	}

	ReqCmd.AddCommand(ReqLsCmd, ReqGetCmd)

	quiet = ReqCmd.PersistentFlags().BoolP("quiet", "q", false, "Return id only")
	noPager = ReqCmd.PersistentFlags().BoolP("no-pager", "P", false, "Disable pager")
	asJson = ReqLsCmd.PersistentFlags().BoolP("as-json", "j", false, "As json")
	asJsonList = ReqLsCmd.PersistentFlags().BoolP("as-json-list", "J", false, "As json list")
	cols = ReqLsCmd.PersistentFlags().StringSliceP("cols", "c", []string{"ID", "Method", "Name", "Modified", "ParentID", "IsPrivate"}, "Set columns")

	rootCmd.AddCommand(ReqCmd)
}
