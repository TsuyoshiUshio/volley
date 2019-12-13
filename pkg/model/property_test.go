package model

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/stretchr/testify/assert"
)

func TestJMeterProperty_GenerateModifiedProperty(t *testing.T) {

	sourcePath := "test-data/jmeter-property/jmeter.properties"
	destinationPath := "jmeter.properties"
	jp := &JMeterProperty{
		RemoteHostIPs: []string{"10.0.0.1", "10.0.0.2"},
	}
	jp.GenerateModifiedProperty(sourcePath, destinationPath)
	source, err := readFile(sourcePath)
	if err != nil {
		panic(err)
	}
	destination, err := readFile(destinationPath)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "-remote_hosts=127.0.0.1\n+remote_hosts=10.0.0.1,10.0.0.2\n", diffOnly(source, destination))

	err = os.Remove(destinationPath)
	if err != nil {
		panic(err)
	}

}

func diffOnly(source, destination string) string {
	r := regexp.MustCompile(`^(\+|\-).*$`)
	diffs := diff.LineDiffAsLines(source, destination)
	var d strings.Builder
	for _, l := range diffs {
		if r.MatchString(l) {
			d.WriteString(l)
			d.WriteString("\n")
		}
	}
	return d.String()
}

func readFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	return string(b), nil
}
