package command

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/urfave/cli/v2"
)

type ConfigCommand struct {
}

func (s *ConfigCommand) Upload(c *cli.Context) error {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fieldName := "file"
	matchJmx := regexp.MustCompile(`.*\.jmx`)
	matchCsv := regexp.MustCompile(`.*\.csv`)
	err := filepath.Walk(c.String("directory"), func(path string, info os.FileInfo, err error) error {
		if matchJmx.MatchString(path) || matchCsv.MatchString(path) {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			fw, err := mw.CreateFormFile(fieldName, path)
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, file)
			if err != nil {
				return err
			}
			file.Close()
		}
		return nil
	})
	if err != nil {
		return err
	}
	contentType := mw.FormDataContentType()
	err = mw.Close()
	if err != nil {
		return err
	}
	resp, err := http.Post(c.String("master")+":"+c.String("port")+"/config", contentType, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(responseBody))
	return nil

}
