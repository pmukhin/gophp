package main

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
	"github.com/pmukhin/gophp/repl"
	"errors"
	"github.com/pmukhin/gophp/interpret"
)

func main() {
	root := &cobra.Command{
		Use:   "gophp",
		Short: "Dialect of PHP written in go",
		Run: func(cmd *cobra.Command, args []string) {
			e := interpret.Main(args[0])
			if e != nil {
				cmd.Printf("an error occured: %v\n", e)
			}
		},
		// if it's just one file ... we run it
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires exactly one arg")
			}
			return nil
		},
	}

	root.AddCommand(&cobra.Command{
		Use:   "repl",
		Short: "A simple REPL",
		Run: func(cmd *cobra.Command, args []string) {
			repl.Main()
		},
	})

	err := root.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
