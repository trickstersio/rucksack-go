package runner

type Args []string

func (args Args) Add(values ...string) Args {
	return Args(append(args, values...))
}
