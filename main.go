package main

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
	"github.com/pmukhin/gophp/repl"
)

func main() {
	root := &cobra.Command{
		Use:   "gophp",
		Short: "dialect of php written in go",
	}
	root.AddCommand(&cobra.Command{
		Use:   "repl",
		Short: "a simple REPL",
		Run: func(cmd *cobra.Command, args []string) {
			repl.Main()
		},
	})

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
