package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/trickstersio/rucksack-go/cmd/rucksack/root"
	"github.com/trickstersio/rucksack-go/cmd/rucksack/secrets"
)

func main() {
	cli := &cobra.Command{
		Use:   "rucksack",
		Short: "Command line tools for rucksack project",
	}

	cli.PersistentFlags().String("name", "", "Application name")
	cli.PersistentFlags().String("env", "development", "Application environment")

	cli.AddCommand(root.NewUp())
	cli.AddCommand(root.NewRun())
	cli.AddCommand(root.NewDown())
	cli.AddCommand(root.NewSeed())

	cli.AddCommand(secrets.New())

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
