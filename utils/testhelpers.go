package utils

import (
	"testing"

	"github.com/nmrshll/winter-is-coming-backend/utils/errors"
)

// WrapFatal is a test helper to make a test fail with a wrapped error
func WrapFatal(t *testing.T, err error, format string, args ...interface{}) {
	t.Fatal(errors.Wrap(err, format, args...))
}
