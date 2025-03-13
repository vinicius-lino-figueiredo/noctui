package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/olekukonko/tablewriter"
)

type Table []map[string]string

func (t Table) Cols() []string {
	c := make(map[string]struct{})
	for _, record := range t {
		for field := range record {
			c[field] = struct{}{}
		}
	}
	var cols []string
	for col := range c {
		cols = append(cols, col)
	}
	return cols
}

func (t Table) Values(index int, cols []string) []string {
	s := make([]string, len(cols))
	for n, col := range cols {
		s[n] = t[index][col]
	}
	return s
}

func (t Table) Render(r io.Writer, cols ...string) error {
	tw := tablewriter.NewWriter(r)
	tw.SetHeader(cols)
	for i := range len(t) {
		tw.Append(t.Values(i, cols))
	}
	tw.Render()
	return nil
}

func Output(r io.Reader, paginate bool) error {
	if paginate {
		return Less(r)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(b)
	return err
}

func Less(r io.Reader) error {
	cmd := exec.Command("/usr/bin/less")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
