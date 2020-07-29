package secrets

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trickstersio/rucksack-go/secrets"
)

func NewEncrypt() *cobra.Command {
	var flags struct {
		Key   string
		Nonce string
		File  string
	}

	command := &cobra.Command{
		Use:   "encrypt [FILE]",
		Short: "Encrypts file using given key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := ReadFile(args[0])

			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", args[0], err)
			}

			encryptor, err := secrets.NewAES(flags.Key, flags.Nonce)

			if err != nil {
				return fmt.Errorf("failed to create encryptor for the file %s: %w", args[0], err)
			}

			data, err := encryptor.Encrypt(raw)

			if err != nil {
				return fmt.Errorf("failed to encrypt file %s: %w", args[0], err)
			}

			if err := WriteFile(flags.File, data); err != nil {
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
