package util

import (
	"bufio"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
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
			continue
		}

		body := resp.Body
		defer body.Close()

		fullFilepath := imgDir + "/" + filename
		file, err := os.Create(fullFilepath)
		if err != nil {
			log.Printf("failed to open file %s: %s\n", fullFilepath, err.Error())
			continue
		}
		defer file.Close()

		if _, err := io.Copy(file, body); err != nil {
			log.Printf("failed to write file %s: %s\n", fullFilepath, err.Error())
		}
		if err = file.Sync(); err != nil {
			panic(err)
		}
	}
}
