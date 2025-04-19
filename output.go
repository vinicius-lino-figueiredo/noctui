package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/olekukonko/tablewriter"
)

// Table represents a slice of rows, where each row is a map of column values.
// Each map key is treated as a column name, and the value is its content.
type Table []map[string]string

// Cols returns a slice containing all unique column names found in the table.
// The order is not guaranteed and depends on map iteration.
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

// Values returns the values of a specific row in the order of given columns.
// Missing keys will result in empty strings in the returned slice.
func (t Table) Values(index int, cols []string) []string {
	s := make([]string, len(cols))
	for n, col := range cols {
		s[n] = t[index][col]
	}
	return s
}

// Render prints the table to the given writer using the specified columns. It
// uses tablewriter to create a formatted ASCII table.
func (t Table) Render(r io.Writer, cols ...string) error {
	tw := tablewriter.NewWriter(r)
	tw.SetHeader(cols)
	for i := range len(t) {
		tw.Append(t.Values(i, cols))
	}
	tw.Render()
	return nil
}

// Output writes the reader's content to stdout or paginates it via less. If
// paginate is true, the content is piped through the "less" command.
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

// Less opens a "less" subprocess and pipes the reader's content to it. It sets
// stdin, stdout, and stderr to the current process's streams.
func Less(r io.Reader) error {
	cmd := exec.Command("/usr/bin/less")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
