package runner

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Runner struct {
	env     string
	name    string
	workdir WorkDir
}

func NewRunner(env string, name string) (*Runner, error) {
	wd, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	if env == "" {
		env = "development"
	}

	runner := &Runner{
		env:  env,
		name: name,
		workdir: WorkDir{
			root: filepath.Clean(wd),
		},
	}

	return runner, nil
}

func (runner *Runner) Up() error {
	if err := runner.prepare(); err != nil {
		return fmt.Errorf("runner preparations failed: %w", err)
	}

	if err := runner.up("api"); err != nil {
		return fmt.Errorf("failed to get api up: %w", err)
	}

	return nil
}

func (runner *Runner) Down() error {
	if err := runner.local("rm", Args{"-rf", runner.workdir.CacheDir()}); err != nil {
		return fmt.Errorf("failed to delete cache directory")
	}

	if err := runner.local("docker-compose", runner.defaultDockerComposeArgs().Add("down", "--volumes")); err != nil {
		return fmt.Errorf("failed to shut environment down: %w", err)
	}

	return nil
}

func (runner *Runner) Seed() error {
	if err := runner.prepare(); err != nil {
		return fmt.Errorf("runner preparations failed: %w", err)
	}

	return runner.seed()
}

func (runner *Runner) Run(name string, args Args) error {
	if err := runner.prepare(); err != nil {
		return fmt.Errorf("runner preparations failed: %w", err)
	}

	return runner.remote(name, args)
}

func (runner *Runner) prepare() error {
	if err := NewDir(runner.workdir.TmpDir()).Create(); err != nil {
		return fmt.Errorf("failed to create tmp folder: %w", err)
	}

	if err := NewDir(runner.workdir.CacheDir()).Create(); err != nil {
		return fmt.Errorf("failed to create cache folder: %w", err)
	}

	if err := runner.up("runner"); err != nil {
		return fmt.Errorf("failed to get runner container up: %w", err)
	}

	if err := runner.cached(runner.modulesCacheKey, runner.installModules); err != nil {
		return fmt.Errorf("failed to install modules: %w", err)
	}

	if err := runner.cached(runner.migrationsCacheKey, runner.migrate); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := runner.cached(runner.seedsCacheKey, runner.seed); err != nil {
		return fmt.Errorf("failed to run seeds: %w", err)
	}

	return nil
}

func (runner *Runner) seed() error {
	if runner.env == "test" {
		return nil
	}

	dir := Dir{
		path: filepath.Join(runner.workdir.root, "db/seed"),
	}

	args := Args{"run"}

	err := dir.EachSourceFile(func(fileInfo os.FileInfo) error {
		args = append(args, filepath.Join("db/seed", fileInfo.Name()))
		return nil
	})

	if err != nil {
		return err
	}

	return runner.remote("go", append(args, "migrate"))

}

func (runner *Runner) migrate() error {
	dir := Dir{
		path: filepath.Join(runner.workdir.root, "db/migrations"),
	}

	args := Args{"run"}

	err := dir.EachSourceFile(func(fileInfo os.FileInfo) error {
		path := filepath.Join("db/migrations", fileInfo.Name())
		log.Println(path)
		args = append(args, path)
		return nil
	})

	if err != nil {
		return err
	}

	return runner.remote("go", append(args, "migrate"))
}

func (runner *Runner) installModules() error {
	if err := runner.remote("go", Args{"mod", "download"}); err != nil {
		return fmt.Errorf("failed to download modules: %w", err)
	}

	if err := runner.remote("go", Args{"install", "github.com/trickstersio/rucksack-go/..."}); err != nil {
		return fmt.Errorf("failed to install rucksack: %w", err)
	}

	return nil
}

func (runner *Runner) cached(keygen func() (string, error), fallback func() error) error {
	key, err := keygen()

	if err != nil {
		return fmt.Errorf("failed to generate cache key: %w", err)
	}

	path := filepath.Join(
		runner.workdir.CacheDir(),
		fmt.Sprintf("%x", sha256.Sum256([]byte(key))),
	)

	file := File{
		path: path,
	}

	if file.Exists() {
		return nil
	}

	if err := fallback(); err != nil {
		return err
	}

	if err := file.Touch(); err != nil {
		return fmt.Errorf("failed to create cache key: %w", err)
	}

	return nil
}

func (runner *Runner) up(service string) error {
	return runner.local("docker-compose", runner.defaultDockerComposeArgs().Add("up", "--detach", service))
}

func (runner *Runner) local(name string, args Args) error {
	cmd := exec.Command(name, args...)

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("APP_ENV=%s", runner.env),
		fmt.Sprintf("APP_NAME=%s", runner.name),
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	log.Println("Running", cmd.String())

	return cmd.Run()
}

func (runner *Runner) remote(name string, args Args) error {
	return runner.local("docker-compose", runner.defaultDockerComposeArgs().Add("exec").Add(
		runner.defaultDockerComposeExecArgs()...).Add("runner").Add(name).Add(args...),
	)
}

func (runner *Runner) defaultDockerComposeArgs() Args {
	return Args{
		fmt.Sprintf("--project-name=%s_%s", runner.name, runner.env),
		fmt.Sprintf("--project-directory=%s", runner.workdir.root),
		fmt.Sprintf(`--file=%s`, runner.workdir.DockerComposeConfig(runner.env)),
	}
}

func (runner *Runner) defaultDockerComposeExecArgs() Args {
	result := Args{}

	if runner.isCI() {
		result = result.Add("-T")
	}

	return result
}

func (runner *Runner) isCI() bool {
	val, ok := os.LookupEnv("GITHUB_ACTIONS")

	if !ok {
		return false
	}

	return val == "true"
}

func (runner *Runner) seedsCacheKey() (string, error) {
	dir := Dir{
		path: filepath.Join(runner.workdir.root, "db/seed"),
	}

	return dir.Digest()
}

func (runner *Runner) migrationsCacheKey() (string, error) {
	dir := Dir{
		path: filepath.Join(runner.workdir.root, "db/migrations"),
	}

	return dir.Digest()
}

func (runner *Runner) modulesCacheKey() (string, error) {
	mod, err := File{
		path: filepath.Join(runner.workdir.root, "go.mod"),
	}.Digest()

	if err != nil {
		return "", fmt.Errorf("failed to calculate mod.go digest: %w", err)
	}

	sum, err := File{
		path: filepath.Join(runner.workdir.root, "go.sum"),
	}.Digest()

	if err != nil {
		return "", fmt.Errorf("failed to calculate go.sum digest: %w", err)
	}

	return mod + sum, nil
}
