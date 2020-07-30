package secrets

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trickstersio/rucksack-go/crypto"
)

func NewEncrypt() *cobra.Command {
	var flags struct {
		Key    string
		Nonce  string
		File   string
		Config string
	}

	var config crypto.Config

	command := &cobra.Command{
		Use:   "encrypt [FILE]",
		Short: "Encrypts file using given key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ReadConfigFile(flags.Config, &config); err == nil {
				if len(flags.Key) == 0 {
					flags.Key = config.Key
				}

				if len(flags.Nonce) == 0 {
					flags.Nonce = config.Nonce
				}
			}

			inputFileName := args[0]
			outputFileName := args[0]

			if len(flags.File) > 0 {
				outputFileName = flags.File
			}

			raw, err := ReadFile(inputFileName)

			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", inputFileName, err)
			}

			encryptor, err := crypto.NewAES(flags.Key, flags.Nonce)

			if err != nil {
				return fmt.Errorf("failed to create encryptor for the file %s: %w", inputFileName, err)
			}

			data, err := encryptor.Encrypt(raw)

			if err != nil {
				return fmt.Errorf("failed to encrypt file %s: %w", inputFileName, err)
			}

			if err := WriteFile(outputFileName, data); err != nil {
				return fmt.Errorf("failed to create output file %s: %w", outputFileName, err)
			}

			return nil
		},
	}

	command.Flags().StringVar(&flags.Key, "key", "", "Encryption key encoded in Base64")
	command.Flags().StringVar(&flags.Nonce, "nonce", "", "Encryption nonce (12 symbols)")
	command.Flags().StringVar(&flags.File, "file", "", "Path to output file")
	command.Flags().StringVar(&flags.Config, "config", "config/key.json", "Path to config file")

	return command
}
