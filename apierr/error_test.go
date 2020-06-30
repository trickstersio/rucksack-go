package apierr

import (
	"testing"

	"github.com/matryer/is"
)

func TestError_Merge(t *testing.T) {
	assert := is.New(t)

	err := New().Add(
		Field("email", New("Email can not be empty")),
	)

	assert.Equal(err.Fields["email"].Messages, []string{
		"Email can not be empty",
	})

	err.Add(
		Field("email", New("Email must be unique")),
	)

	assert.Equal(err.Fields["email"].Messages, []string{
		"Email can not be empty",
		"Email must be unique",
	})

	err.Add(
		Field("tasks", Array(
			0,
			Field("name", New("Task name can not be empty"))),
		),
	)

	assert.Equal(err.Fields["tasks"].Fields["0"].Fields["name"].Messages, []string{
		"Task name can not be empty",
	})
}
