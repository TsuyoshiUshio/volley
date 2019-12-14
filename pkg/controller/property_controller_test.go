package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/TsuyoshiUshio/volley/pkg/helper"
	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJUpdateJMeterConfig_NormalCase(t *testing.T) {
	inputJson := "{\"remote_host_ips\": [\"10.0.0.1\",\"10.0.0.2\"]}"
	expectedJMeterPropertyFilePath := filepath.Join("foo", "bar", "jmeter.properties")

	actualResponseWriter := &MyResponseWriter{
		WrittenResponse: strings.Builder{},
		WrittenHeaders:  []int{},
	}
	context := gin.Context{
		Request: &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(inputJson))),
		},
		Writer: actualResponseWriter,
	}

	monkey.Patch(helper.GetJMeterPropertyFilePath, func() string {
		return expectedJMeterPropertyFilePath
	})
	monkey.Patch(os.Stat, func(file string) (os.FileInfo, error) {
		assert.Equal(t, model.JMeterPropertyFile, file)
		return nil, fmt.Errorf("File not found %s", file)
	})
	var p *model.JMeterProperty

	monkey.PatchInstanceMethod(reflect.TypeOf(p), "GenerateModifiedProperty", func(p *model.JMeterProperty, source, destination string) error {
		assert.Equal(t, expectedJMeterPropertyFilePath, source)
		return nil
	})
	defer monkey.UnpatchAll()
	UpdateJMeterConfig(&context)

	assert.Equal(t, "{\"remote_host_ips\":[\"10.0.0.1\",\"10.0.0.2\"]}\n", actualResponseWriter.WrittenResponse.String())
	assert.Equal(t, actualResponseWriter.WrittenHeaders[0], 200)
}

func TestJUpdateJMeterConfig_Conifg_Already_Exists(t *testing.T) {
	inputJson := "{\"remote_host_ips\": [\"10.0.0.1\",\"10.0.0.2\"]}"
	expectedJMeterPropertyFilePath := filepath.Join("foo", "bar", "jmeter.properties")

	actualResponseWriter := &MyResponseWriter{
		WrittenResponse: strings.Builder{},
		WrittenHeaders:  []int{},
	}
	context := gin.Context{
		Request: &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(inputJson))),
		},
		Writer: actualResponseWriter,
	}

	monkey.Patch(helper.GetJMeterPropertyFilePath, func() string {
		return expectedJMeterPropertyFilePath
	})
	monkey.Patch(os.Stat, func(file string) (os.FileInfo, error) {
		assert.Equal(t, model.JMeterPropertyFile, file)
		return nil, nil
	})
	monkey.Patch(os.Remove, func(file string) error {
		assert.Equal(t, model.JMeterPropertyFile, file)
		return nil
	})

	var p *model.JMeterProperty

	monkey.PatchInstanceMethod(reflect.TypeOf(p), "GenerateModifiedProperty", func(p *model.JMeterProperty, source, destination string) error {
		assert.Equal(t, expectedJMeterPropertyFilePath, source)
		return nil
	})
	defer monkey.UnpatchAll()
	UpdateJMeterConfig(&context)

	assert.Equal(t, "{\"remote_host_ips\":[\"10.0.0.1\",\"10.0.0.2\"]}\n", actualResponseWriter.WrittenResponse.String())
	assert.Equal(t, actualResponseWriter.WrittenHeaders[0], 200)
}

type MyResponseWriter struct {
	WrittenResponse strings.Builder
	WrittenHeaders  []int
	gin.ResponseWriter
}

func (w *MyResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *MyResponseWriter) Write(p []byte) (n int, err error) {
	w.WrittenResponse.WriteString(string(p))
	return len(p), nil
}

func (w *MyResponseWriter) WriteHeader(code int) {
	w.WrittenHeaders = append(w.WrittenHeaders, code)
}
