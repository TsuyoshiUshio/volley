package helper

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZipScenario(t *testing.T) {
	// Zip one directry recursively
	// extract zip file in a temp folder
	// check if there is a file on the deepest directry
	tempPath := filepath.Join(".", ".test")

	// If the tempPath exists, remove it
	if _, err := os.Stat(tempPath); err == nil {
		os.RemoveAll(tempPath)
	}

	os.MkdirAll(tempPath, os.ModePerm)
	zipFilePath := filepath.Join(tempPath, "hello.zip")
	Zip(filepath.Join("test-data", "zip", "hello"), zipFilePath)
	UnZip(zipFilePath, tempPath)
	textPath := filepath.Join(tempPath, "hello", "world.txt")
	if _, err := os.Stat(textPath); os.IsNotExist(err) {
		assert.Fail(t, "Unzipped file  manifest can not found at :"+textPath)
	}

	// clean up

	os.RemoveAll(tempPath)
}
