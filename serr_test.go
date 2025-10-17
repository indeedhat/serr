package serr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedWrapText = `Error: base error
Context: {"multi_1":["multi","one"],"multi_2":["multi","two"],"single":["single data"]}`

func TestWrap(t *testing.T) {
	err := errors.New("base error")
	se := Wrap(err)
	se.AddContext("single", "single data")
	se.AddContextMap(Context{
		"multi_1": {"multi", "one"},
		"multi_2": {"multi", "two"},
	})

	require.Equal(t, Context{
		"single":  {"single data"},
		"multi_1": {"multi", "one"},
		"multi_2": {"multi", "two"},
	}, se.Context())

	require.Equal(t, expectedWrapText, se.Error())
}

var expectedContextWrapText = `Error: base error
Context: {"base":["data"],"multi_1":["multi","one"],"multi_2":["multi","two"],"single":["single data"]}`

func TestContextWrap(t *testing.T) {
	ectx := Context{"base": {"data"}}
	err := errors.New("base error")
	se := ectx.Wrap(err)
	se.AddContext("single", "single data")
	se.AddContextMap(Context{
		"multi_1": {"multi", "one"},
		"multi_2": {"multi", "two"},
	})

	require.Equal(t, Context{
		"base":    {"data"},
		"single":  {"single data"},
		"multi_1": {"multi", "one"},
		"multi_2": {"multi", "two"},
	}, se.Context())

	require.Equal(t, expectedContextWrapText, se.Error())

	// make sure that subsequent changes do not effect the base context
	require.Equal(t, Context{"base": {"data"}}, ectx)
}
