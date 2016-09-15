package testpack

import (
	"testing"
	"os"
	"reflect"
)

// GetTestName returns test name
func GetTestName(t interface{}) (n string) {
	pointerVal := reflect.ValueOf(t)
	val := reflect.Indirect(pointerVal)

	member := val.FieldByName("name")
	n = member.String()
	return
}

type skipper interface {
	Skip(args ...interface{})
	Skipf(format string, args ...interface{})
}

// SkipUnless skip tests
func SkipUnless(t skipper, short bool, env ...string) {
	if short && testing.Short() {
		t.Skip("wan't run in -short")
		return
	}
	for _, e := range env {
		if os.Getenv(e) == "" {
			t.Skipf("%s OS environment variable is not defined", e)
		}
	}

}
