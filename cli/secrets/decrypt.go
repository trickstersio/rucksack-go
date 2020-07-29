package secrets

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trickstersio/rucksack-go/secrets"
)

func NewDecrypt() *cobra.Command {
	var flags struct {
		Key   string
		Nonce string
		File  string
	}

	command := &cobra.Command{
		Use:   "decrypt [FILE]",
		Short: "Decrypts file using given key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := ReadFile(args[0])

			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", args[0], err)
			}

			decryptor, err := secrets.NewAES(flags.Key, flags.Nonce)

			if err != nil {
				return fmt.Errorf("failed to create decryptor file %s: %w", args[0], err)
			}

			raw, err := decryptor.Decrypt(data)

			if err != nil {
				return fmt.Errorf("failed to decrypt file %s: %w", args[0], err)
			}

			if err := WriteFile(flags.File, raw); err != nil {
				return fmt.Errorf("failed to create output file %s: %w", flags.File, err)
			}

			return nil
		},
	}

	command.Flags().StringVar(&flags.Key, "key", "", "Encryption key encoded in Base64")
	command.Flags().StringVar(&flags.Nonce, "nonce", "", "Encryption nonce (12 symbols)")
	command.Flags().StringVar(&flags.File, "file", "", "Path to output file")

	return command
}
