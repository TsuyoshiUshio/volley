package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// JMeterPropertyDir is the directory that include property.json
	JMeterPropertyDir = "property"
	// JMeterPropertyJSON is the JSON file that is serialized from JMeterProperty instance.
	JMeterPropertyJSON = "property.json"
	// JMeterPropertyFile is the file name of the jmeter.property file
	JMeterPropertyFile = "jmeter.properties"
)

// JMeterProperty is a struct that represent updated data of jmeter.properties
type JMeterProperty struct {
	RemoteHostIPs []string `json:"remote_host_ips" binding:"required"`
}

// NewJMeterProperty will create an instance from json file on property directory
func NewJMeterProperty() (*JMeterProperty, error) {
	jmeterPropertyFilePath := filepath.Join(".", JMeterPropertyDir, JMeterPropertyJSON)
	b, err := ioutil.ReadFile(jmeterPropertyFilePath)
	if err != nil {
		return nil, err
	}
	var jmeterProperty JMeterProperty
	err = json.Unmarshal(b, &jmeterProperty)
	if err != nil {
		return nil, err
	}
	return &jmeterProperty, nil
}

// WriteFile write JSON marshalled JMeterProperty
func (p *JMeterProperty) WriteFile() error {
	b, err := json.Marshal(*p)
	if err != nil {
		return err
	}
	jmeterPropertyDirPath := filepath.Join(".", JMeterPropertyDir)
	if _, err := os.Stat(jmeterPropertyDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(jmeterPropertyDirPath, os.ModePerm)
	}
	jmeterPropertyJSONPath := filepath.Join(jmeterPropertyDirPath, JMeterPropertyJSON)

	err = ioutil.WriteFile(jmeterPropertyJSONPath, b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// CommaSeparatedRemoteHostIPs is string representation of RemoteHostIPs
func (p *JMeterProperty) CommaSeparatedRemoteHostIPs() string {
	return strings.Join(p.RemoteHostIPs, ",")
}
