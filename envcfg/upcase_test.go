package envcfg

import (
	"testing"

	"github.com/matryer/is"
)

func Test_UpcaseLookuper(t *testing.T) {
	assert := is.New(t)

	lookuper := UpcaseLookuper(LookuperFunc(func(key string) (string, bool) {
		assert.Equal(key, "MY_KEY")
		return "MY_VALUE", true
	}))

	value, ok := lookuper.Lookup("my_key")

	assert.True(ok)
	assert.Equal(value, "MY_VALUE")
}
