package envcfg

type LookuperFunc func(key string) (string, bool)

func (f LookuperFunc) Lookup(key string) (string, bool) {
	return f(key)
}
