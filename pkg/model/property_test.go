package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJMeterProperty_WriteFile(t *testing.T) {
	inputJMeterProperty := &JMeterProperty{
		RemoteHostIPs: []string{"10.0.0.4", "10.0.0.5"},
	}
	err := inputJMeterProperty.WriteFile()
	targetFilePath := filepath.Join(".", JMeterPropertyDir, JMeterPropertyJSON)
	assert.Nil(t, err)
	defer os.Remove(targetFilePath)
	file, err := ioutil.ReadFile(targetFilePath)
	assert.Nil(t, err)
	var actualJMeterProperty JMeterProperty
	json.Unmarshal(file, &actualJMeterProperty)
	assert.Equal(t, inputJMeterProperty.RemoteHostIPs[0], actualJMeterProperty.RemoteHostIPs[0])
	assert.Equal(t, inputJMeterProperty.RemoteHostIPs[1], actualJMeterProperty.RemoteHostIPs[1])
}

func TestJMeterProperty_CommaSeparatedRemoteHostIPs_Multi(t *testing.T) {
	inputJMeterProperty := &JMeterProperty{
		RemoteHostIPs: []string{"10.0.0.4", "10.0.0.5"},
	}
	assert.Equal(t, "10.0.0.4,10.0.0.5", inputJMeterProperty.CommaSeparatedRemoteHostIPs())
}

func TestJMeterProperty_CommaSeparatedRemoteHostIPs_Single(t *testing.T) {
	inputJMeterProperty := &JMeterProperty{
		RemoteHostIPs: []string{"10.0.0.4"},
	}
	assert.Equal(t, "10.0.0.4", inputJMeterProperty.CommaSeparatedRemoteHostIPs())
}

func TestJMeterProperty_CommaSeparatedRemoteHostIPs_Zero(t *testing.T) {
	inputJMeterProperty := &JMeterProperty{
		RemoteHostIPs: []string{},
	}
	assert.Equal(t, "", inputJMeterProperty.CommaSeparatedRemoteHostIPs())
}
