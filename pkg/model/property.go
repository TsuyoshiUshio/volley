package model

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type JMeterProperty struct {
	RemoteHostIPs []string `json:"remote_host_ips" binding:"required"`
}

func (p *JMeterProperty) GenerateModifiedProperty(source, destination string) error {
	// read from source
	// write to the destination
	// read the source line by line
	// if find the match, replace it
	s, err := os.Open(source)
	if err != nil {
		return err
	}
	defer s.Close()

	if _, err := os.Stat(destination); err == nil {
		err = os.Remove(destination)
		if err != nil {
			return err
		}
	}

	d, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer d.Close()

	scanner := bufio.NewScanner(s)
	r := regexp.MustCompile(`^remote_hosts=.*$`)
	for scanner.Scan() {
		line := scanner.Text()
		if r.MatchString(line) {
			fmt.Fprintln(d, "remote_hosts="+strings.Join(p.RemoteHostIPs, ","))
		} else {
			fmt.Fprintln(d, line)
		}
	}

	return nil
}
