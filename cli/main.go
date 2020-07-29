package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/trickstersio/rucksack-go/cli/secrets"
)

func main() {
	root := &cobra.Command{
		Use:   "rucksack",
		Short: "Command line tools for rucksack project",
	}

	root.AddCommand(secrets.New())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
