package testpack

import (
	"testing"
	"os"
	"reflect"
	"regexp"
)

// GetTestName returns test name with
func GetTestName(t interface{}) (n string) {
	pointerVal := reflect.ValueOf(t)
	val := reflect.Indirect(pointerVal)

	member := val.FieldByName("name")
	n = member.String()
	return
}

func GetTestNameN(t interface{}) (n string) {
	n = NormalizeName(GetTestName(t))
	return
}

func NormalizeName(in string) (out string)  {
	re := regexp.MustCompile("[^a-zA-Z0-9_]")
	out = re.ReplaceAllLiteralString(in, "_")
	return
}

type skipper interface {
	Skip(args ...interface{})
	Skipf(format string, args ...interface{})
}

// SkipUnless skip tests
func SkipUnless(t skipper, short bool, env ...string) {
	if short && testing.Short() {
		t.Skipf("%s: wan't run with -short", GetTestName(t))
		return
	}
	for _, e := range env {
		if os.Getenv(e) == "" {
			t.Skipf("%s OS environment variable is not defined", e)
		}
	}

}
