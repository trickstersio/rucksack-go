package secrets

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trickstersio/rucksack-go/secrets"
)

func NewTokens() *cobra.Command {
	command := &cobra.Command{
		Use:   "tokens",
		Short: "Create and validate secure tokens",
	}

	command.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Creates new secure token",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := secrets.RandomToken()

			if err != nil {
				return fmt.Errorf("failed to generate random token: %w", err)
			}

			digest, err := token.Digest()

			if err != nil {
				return fmt.Errorf("failed to generate token digest: %w", err)
			}

			fmt.Printf("Token: %s\n", token.String())
			fmt.Printf("Digest: %s\n", digest)

			return nil
		},
	})

	command.AddCommand(&cobra.Command{
		Use:   "validate TOKEN DIGEST",
		Short: "Validates provided token against digest",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			token := secrets.NewToken(args[0])

			if err := token.Validate(args[1]); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			return nil
		},
	})

	return command
}
