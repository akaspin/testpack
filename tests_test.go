package testpack_test

import (
	"testing"
	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/akaspin/testpack"
)

func TestGetTestName(t *testing.T) {
	assert.Equal(t, testpack.GetTestName(t), "TestGetTestName")

	t.Run("strange=test name", func(t *testing.T) {
		assert.Equal(t, testpack.GetTestName(t),
			"TestGetTestName/strange=test_name") // note about underscore
	})
}
