package runner

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
)

type File struct {
	path string
}

func NewFile(path string) File {
	return File{
		path: path,
	}
}

func (f File) Exists() bool {
	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		return false
	}

	return true
}

func (f File) Touch() error {
	if _, err := os.Create(f.path); err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

func (f File) Digest() (string, error) {
	data, err := ioutil.ReadFile(f.path)

	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}
