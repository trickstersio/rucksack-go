package secrets

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "secrets",
		Short: "Secrets management: create, edit, encrypt, decrypt",
	}

	command.AddCommand(NewEncrypt())
	command.AddCommand(NewDecrypt())

	return command
}

func WriteFile(filename string, data []byte) error {
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func ReadFile(filename string) ([]byte, error) {
	fileInfo, err := os.Stat(filename)

	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("file can not be directory")
	}

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf("failed to read content from: %w", err)
	}

	return data, nil
}
