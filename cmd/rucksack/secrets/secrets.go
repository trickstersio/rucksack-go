package secrets

import (
	"encoding/json"
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

func CheckFileExists(filename string) error {
	fileInfo, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %w", err)
		}

		return fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("file can not be directory")
	}

	return nil
}

func ReadConfigFile(filename string, out interface{}) error {
	if err := CheckFileExists(filename); err != nil {
		return nil
	}

	data, err := ReadFile(filename)

	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return nil
}

func WriteFile(filename string, data []byte) error {
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func ReadFile(filename string) ([]byte, error) {
	err := CheckFileExists(filename)

	if err != nil {
		return nil, fmt.Errorf("failed to check that file exists: %w", err)
	}

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf("failed to read content from: %w", err)
	}

	return data, nil
}
