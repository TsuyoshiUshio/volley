package helper

import (
	"github.com/TsuyoshiUshio/volley/pkg/model"
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
	return filepath.Join(jmeterBinDir, model.JMeterPropertyFile)
}
