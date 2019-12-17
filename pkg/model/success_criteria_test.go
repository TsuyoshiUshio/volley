package model

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessCriteria_Validate_NormalCase(t *testing.T) {
	fileName := filepath.Join("test-data", "success-criteria", "avg-time-error-on-rps", "stress.log")
	// This fixture will cause this.
	// 2019/12/17 00:29:11 TotalRequest: 3916, Average Latency: 11593, ErrorRatio: 1 %
	// 2019/12/17 00:29:11 Request Per Second Up to: 250, Target Average Letency Less than: 0, Target Error Ratio Less than: 10

	m := make(map[string]int64)
	m[REQUEST_PER_SECOND] = 250
	m[AVERAGE_LATENCY] = 20000
	m[ERROR_RATIO] = 10
	criteria := &AverageTimeAndErrorOnRPSSuccessCriteria{
		SuccessCriteria: SuccessCriteria{
			Name:       "average_time_error_on_rps",
			Parameters: m,
		},
	}
	isSuccess, err := criteria.Validate(fileName)
	if err != nil {
		panic(err)
	}
	assert.True(t, isSuccess)
}

func TestSuccessCriteria_Validate_ExceedAverageLatency(t *testing.T) {
	fileName := filepath.Join("test-data", "success-criteria", "avg-time-error-on-rps", "stress.log")
	// This fixture will cause this.
	// 2019/12/17 00:29:11 TotalRequest: 3916, Average Latency: 11593, ErrorRatio: 1 %
	// 2019/12/17 00:29:11 Request Per Second Up to: 250, Target Average Letency Less than: 0, Target Error Ratio Less than: 10

	m := make(map[string]int64)
	m[REQUEST_PER_SECOND] = 250
	m[AVERAGE_LATENCY] = 10000
	m[ERROR_RATIO] = 10
	criteria := &AverageTimeAndErrorOnRPSSuccessCriteria{
		SuccessCriteria: SuccessCriteria{
			Name:       "average_time_error_on_rps",
			Parameters: m,
		},
	}
	isSuccess, err := criteria.Validate(fileName)
	if err != nil {
		panic(err)
	}
	assert.False(t, isSuccess)
}

func TestSuccessCriteria_Validate_ExceedErrorRatio(t *testing.T) {
	fileName := filepath.Join("test-data", "success-criteria", "avg-time-error-on-rps", "stress.log")
	// This fixture will cause this.
	// 2019/12/17 00:29:11 TotalRequest: 3916, Average Latency: 11593, ErrorRatio: 1 %
	// 2019/12/17 00:29:11 Request Per Second Up to: 250, Target Average Letency Less than: 0, Target Error Ratio Less than: 10

	m := make(map[string]int64)
	m[REQUEST_PER_SECOND] = 250
	m[AVERAGE_LATENCY] = 20000
	m[ERROR_RATIO] = 0
	criteria := &AverageTimeAndErrorOnRPSSuccessCriteria{
		SuccessCriteria: SuccessCriteria{
			Name:       "average_time_error_on_rps",
			Parameters: m,
		},
	}
	isSuccess, err := criteria.Validate(fileName)
	if err != nil {
		panic(err)
	}
	assert.False(t, isSuccess)
}

func TestSuccessCriteria_Instanciate(t *testing.T) {
	configFilePath := filepath.Join("test-data", "success-criteria", "config", "success_criteria.json")
	criteria, err := NewAverageTimeAndErrorOnRPSSuccessCriteria(configFilePath)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "average_time_error_on_rps", criteria.Name)
	assert.Equal(t, int64(10000), criteria.Parameters[AVERAGE_LATENCY])
	assert.Equal(t, int64(10), criteria.Parameters[ERROR_RATIO])
	assert.Equal(t, int64(250), criteria.Parameters[REQUEST_PER_SECOND])
}
