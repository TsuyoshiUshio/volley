package model

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJMeterLog_NormalCase(t *testing.T) {
	jMeterLog := &JMeterLog{}
	jMeterLog.InitializeWithFile(filepath.Join("test-data", "success-criteria", "avg-time-error-on-rps", "stress.log"))
	// The header is
	// timeStamp,elapsed,label,responseCode,responseMessage,threadName,dataType,success,failureMessage,bytes,sentBytes,grpThreads,allThreads,URL,Latency,IdleTime,Connect
	inputLine := "1575859346249,5173,addNew - Success,200,OK,Ultimate Thread Group - addNew 1-92,text,true,,1014,298,100,100,https://localhost:4500/api/v1/AddNew,1269,0,2"
	values, _ := jMeterLog.Parse(inputLine)
	assert.Equal(t, "1575859346249", (*values)["timeStamp"])
	assert.Equal(t, "2", (*values)["Connect"])

}
