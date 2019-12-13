package helper

import (
	"log"
	"os/exec"
	"path/filepath"
)

func GetJMeterPropertyFilePath() string {
	jmeterPath, err := exec.LookPath("jmeter")
	if err != nil {
		log.Fatal("Can not find jmeter executable on the path enviornment variables.", err)
	}
	jmeterBinDir := filepath.Dir(jmeterPath)
	return jmeterBinDir
}
