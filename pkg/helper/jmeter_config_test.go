package helper

import (
	"os/exec"
	"path/filepath"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestGetJMeterPropertyFilePath(t *testing.T) {
	var expectedExecutable = "jmeter"
	var expectedExecutablePath = filepath.Join("foo", "bar", expectedExecutable)
	var expectedPath = filepath.Join("foo", "bar")
	monkey.Patch(exec.LookPath, func(file string) (string, error) {
		assert.Equal(t, expectedExecutable, file)
		return expectedExecutablePath, nil
	})
	defer monkey.UnpatchAll()
	actualPath := GetJMeterPropertyFilePath()
	assert.Equal(t, expectedPath, actualPath)
}
