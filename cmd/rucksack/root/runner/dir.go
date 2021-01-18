package runner

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Dir struct {
	path string
}

func NewDir(path string) Dir {
	return Dir{
		path: path,
	}
}

func (dir Dir) Create() error {
	if err := os.MkdirAll(dir.path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir.path, err)
	}

	return nil
}

func (dir Dir) EachSourceFile(do func(fileInfo os.FileInfo) error) error {
	entries, err := ioutil.ReadDir(dir.path)

	if err != nil {
		return fmt.Errorf("failed to list entries in directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
			continue
		}

		if err := do(entry); err != nil {
			return err
		}
	}

	return nil
}

func (dir Dir) Digest() (string, error) {
	var digests []string

	err := dir.EachSourceFile(func(fileInfo os.FileInfo) error {
		path := filepath.Join(dir.path, fileInfo.Name())
		digest, err := NewFile(path).Digest()

		if err != nil {
			return fmt.Errorf("failed to calculate digest for %s: %w", path, err)
		}

		digests = append(digests, digest)

		return nil
	})

	if err != nil {
		return "", err
	}

	return strings.Join(digests, ""), nil
}
