package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	// REQUEST_PER_SECOND is Request per second
	REQUEST_PER_SECOND = "rps"
	// AVERAGE_LATENCY is Average Latency
	AVERAGE_LATENCY = "avg_latency"
	// ERROR_RATIO is Error ratio
	ERROR_RATIO = "error_ratio"
)

// ISuccessCriteria is an interface for SuccessCriteria
type ISuccessCriteria interface {
	Validate(fileName string) (bool, error)
}

// SuccessCriteria is a struct that represents Scucess Criteria config
type SuccessCriteria struct {
	Name       string `json:"criteria" binding:"required"`
	Parameters map[string]int64
}

// NewAverageTimeAndErrorOnRPSSuccessCriteria constructs AverageTimeAndErrorOnRPSSuccessCriteria instance from config file path
func NewAverageTimeAndErrorOnRPSSuccessCriteria(configFilePath string) (*AverageTimeAndErrorOnRPSSuccessCriteria, error) {
	config, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	var criteria AverageTimeAndErrorOnRPSSuccessCriteria
	err = json.Unmarshal(config, &criteria)
	if err != nil {
		return nil, fmt.Errorf("Can not parse sucess criteria config fileName: %s, contents: %s, error:%v", configFilePath, string(config), err)
	}
	return &criteria, nil
}

// AverageTimeAndErrorOnRPSSuccessCriteria is a success criteria of Stress Testing.
// This criteria represents, validate the test result with Averate Latency and Error ratio and up to the target Request per second
type AverageTimeAndErrorOnRPSSuccessCriteria struct {
	SuccessCriteria
}

// Validate if the Stress testing is success or not with boolean.
func (c *AverageTimeAndErrorOnRPSSuccessCriteria) Validate(fileName string) (bool, error) {
	// Prepare the parser
	parser := &JMeterLog{}
	parser.InitializeWithFile(fileName)

	// Read the file
	f, err := os.Open(fileName)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Skip the first line. First line is the header
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return false, fmt.Errorf("Cannot read the first line of the JMeter Execution Log: %s", fileName)
	}
	// timeStamp,elapsed,label,responseCode,responseMessage,threadName,dataType,success,failureMessage,bytes,sentBytes,grpThreads,allThreads,URL,Latency,IdleTime,Connect
	requestPerSecond := c.Parameters[REQUEST_PER_SECOND]
	targetAverageLatency := c.Parameters[AVERAGE_LATENCY]
	targetErrorRatio := c.Parameters[ERROR_RATIO]

	var counter int64 = 0
	var lineCounter int64 = 0
	var latencySum int64 = 0
	var errorCounter int64 = 0
	for scanner.Scan() {
		lineString := scanner.Text()
		// if you find / at the end of the lineString, it means it includes multiple line text. Scan it until you find / as the frist string and concat these two.
		lineString, skipCounter, err := convertWhenItIsMultiline(lineString, scanner)
		lineCounter += int64(skipCounter)
		if err != nil {
			return false, fmt.Errorf("Multilined parse error: %s line: %d, error: %v", fileName, lineCounter, err)
		}

		line, err := parser.Parse(lineString)
		if err != nil {
			return false, err
		}
		counter++
		lineCounter++

		threads, err := strconv.ParseInt((*line)["allThreads"], 10, 64)
		if err != nil {
			return false, fmt.Errorf("Can not parse allThreads: %s to int, %s line: %d, error: %v", (*line)["allThreads"], fileName, lineCounter, err)
		}
		if requestPerSecond <= threads {
			isSuccess, err := strconv.ParseBool((*line)["success"])
			if err != nil {
				return false, fmt.Errorf("Can not parse succss: %s to boolean, %s line: %d, error: %v", (*line)["success"], fileName, lineCounter, err)
			}
			if !isSuccess {
				errorCounter++
			}
			latency, err := strconv.ParseInt((*line)["Latency"], 10, 64)
			if err != nil {
				return false, fmt.Errorf("Can not parse Latency: %s to int, %s line: %d, error: %v", (*line)["Latency"], fileName, lineCounter, err)
			}
			latencySum += latency
		}

	}
	averageLatency := latencySum / counter
	errorRatio := (errorCounter * 100) / counter

	log.Printf("TotalRequest: %d, Average Latency: %d, ErrorRatio: %d %%", counter, averageLatency, errorRatio)
	log.Printf("Request Per Second Up to: %d, Target Average Letency Less than: %d, Target Error Ratio Less than: %d", requestPerSecond, targetAverageLatency, targetErrorRatio)

	if averageLatency >= targetAverageLatency || errorRatio >= targetErrorRatio {
		return false, nil
	}

	return true, nil
}

func convertWhenItIsMultiline(line string, scanner *bufio.Scanner) (string, int, error) {
	counter := 0
	var builder strings.Builder
	if strings.HasSuffix(line, "/") {
		builder.WriteString(line)
		for scanner.Scan() {
			newLine := scanner.Text()
			counter++
			if strings.HasPrefix(newLine, "/") {
				builder.WriteString(newLine)
				return builder.String(), counter, nil
			} else {
				builder.WriteString(newLine)
			}
		}
	} else {
		return line, 0, nil
	}
	return builder.String(), counter, fmt.Errorf("It found / at the end of the line on a log file, however, can't find the / that match of the begin of a line. line: %s", builder.String())
}
