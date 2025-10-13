package util

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type ImgFetchYaml struct {
	Root map[string]string `yaml:"root"`
}

func FetchImgFromYaml(imgDir string, imgYamlFilename string) {
	var yamlFile ImgFetchYaml
	open, err := os.Open(imgDir + "/" + imgYamlFilename)
	if err != nil {
		panic(err)
	}

	if err := yaml.NewDecoder(bufio.NewReader(open)).Decode(&yamlFile); err != nil {
		panic(err)
	}

	for filename, resource := range yamlFile.Root {
		resp, err := http.Get(resource)
		if err != nil {
			log.Printf("failed to fetch %s: %s\n", resource, err.Error())
			panic("failed to fetch " + resource)
		}
		if resp.StatusCode != 200 {
			panic("failed to fetch " + resource)
		}

		writeToFile(resp, imgDir, filename)
	}
}

func writeToFile(resp *http.Response, imgDir string, filename string) {
	body := resp.Body
	defer body.Close()

	fullFilepath := imgDir + "/" + filename
	file, err := os.Create(fullFilepath)
	if err != nil {
		panic("failed to open file: " + fullFilepath)
	}
	defer file.Close()

	if _, err := io.Copy(file, body); err != nil {
		panic(fmt.Errorf("failed to write file %s: %v", fullFilepath, err.Error()))
	}
	if err = file.Sync(); err != nil {
		panic(err)
	}
}
