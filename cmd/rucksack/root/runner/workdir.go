package runner

import (
	"fmt"
	"path/filepath"
)

type WorkDir struct {
	root string
	name string
	env  string
}

func (wd WorkDir) CacheDir() string {
	return fmt.Sprintf("%s/runner/%s/%s", wd.TmpDir(), wd.name, wd.env)
}

func (wd WorkDir) TmpDir() string {
	return filepath.Clean(filepath.Join(wd.root, "tmp"))
}

func (wd WorkDir) ConfigDir() string {
	return filepath.Clean(filepath.Join(wd.root, "config"))
}

func (wd WorkDir) DockerComposeConfig(env string) string {
	return filepath.Join(wd.ConfigDir(), fmt.Sprintf("%s/docker-compose.yml", env))
}
