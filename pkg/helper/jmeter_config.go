package helper

import (
	"github.com/TsuyoshiUshio/volley/pkg/model"
	"log"
	"os/exec"
	"path/filepath"
)

// GetJMeterPropertyFilePath returns Jmeter property file path
func GetJMeterPropertyFilePath() string {
	jmeterBinDir := GetJMeterBinDir()
	return filepath.Join(jmeterBinDir, model.JMeterPropertyFile)
}

// GetJMeterBinDir returns path of the JMETER_HOME/bin directory
func GetJMeterBinDir() string {
	jmeterPath, err := exec.LookPath("jmeter")
	if err != nil {
		log.Fatal("Can not find jmeter executable on the path enviornment variables.", err)
	}
	return filepath.Dir(jmeterPath)
}
