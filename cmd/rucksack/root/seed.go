package root

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trickstersio/rucksack-go/cmd/rucksack/root/runner"
)

func NewSeed() *cobra.Command {
	var flags struct {
		Env  string
		Name string
	}

	command := &cobra.Command{
		Use:   "seed",
		Short: "Seeds database with data",
		RunE: func(cmd *cobra.Command, args []string) error {
			runner, err := runner.NewRunner(flags.Env, flags.Name)

			if err != nil {
				return fmt.Errorf("failed to create runner: %w", err)
			}

			if err := runner.Seed(); err != nil {
				return fmt.Errorf("failed to run seed: %w", err)
			}

			return nil
		},
	}

	command.Flags().StringVar(&flags.Env, "env", "development", "Application environment")
	command.Flags().StringVar(&flags.Name, "name", "", "Application name")

	return command
}
