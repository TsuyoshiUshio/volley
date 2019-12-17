package model

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// JMeterLog represent a log file that is generated after JMeter execution.
type JMeterLog struct {
	Header map[int]string
}

// InitializeWithFile parse the first line as a comma separated value.
// Then set the Header by the result of the parse.
func (l *JMeterLog) InitializeWithFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	line := scanner.Text()
	headers := strings.Split(line, ",")
	if len(headers) == 1 {
		log.Printf("Warning: The number of the header of the JMeter log file: %s is 1. Double check if the log has the proper header value.", filePath)
	}
	m := make(map[int]string)
	for i, header := range headers {
		m[i] = header
	}
	l.Header = m
	return nil
}

// Parse the JMeter Logfile line
// The result is the map object that has full line of
func (l *JMeterLog) Parse(line string) (*map[string]string, error) {
	values := strings.Split(line, ",")
	m := make(map[string]string)
	for i, value := range values {
		var header = l.Header[i]
		if header == "" {
			return nil, fmt.Errorf("JMeterLog key value is empty index: %d, key: %s value: %s", i, header, value)
		}
		m[header] = value
	}
	return &m, nil
}
