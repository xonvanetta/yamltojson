package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/icza/dyno"
	"github.com/koding/multiconfig"
	yaml2 "sigs.k8s.io/yaml"
)

var yamlDelim = []byte("---")
var kind = regexp.MustCompile(`.*kind: (.*?)\n`)
var name = regexp.MustCompile(`.*  name: (.*?)\n`)

type config struct {
	StartBufferSize  uint64
	MaxScanTokenSize int
	File             string
}

func isYaml(data []byte) bool {
	if strings.TrimSpace(string(data)) == string(yamlDelim) {
		return false
	}
	var yamlObject interface{}
	return yaml2.Unmarshal(data, &yamlObject) == nil
}

func yamlAdvancedSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {

	firstYamlDelim := bytes.Index(data, yamlDelim)
	if firstYamlDelim == -1 && atEOF == true {
		return 0, data, bufio.ErrFinalToken
	}
	if firstYamlDelim == -1 && atEOF == false {
		return 0, nil, nil
	}

	index := firstYamlDelim + 3
	shouldReturn := false
	if isYaml(data[:index]) {
		shouldReturn = true
		//fmt.Println(string(data[:index]))
		//fmt.Println(string(data))
		//fmt.Println("text")
		//return index, data[:index], nil
	}

	nextYamlDelim := bytes.Index(data[index:], yamlDelim)
	if nextYamlDelim == -1 && shouldReturn && atEOF == true {
		return 0, data[:index], bufio.ErrFinalToken
	}
	if shouldReturn {
		return len(data[:index]), data[:index], nil
	}

	//fmt.Println(index, nextYamlDelim)
	//fmt.Println(atEOF)
	//fmt.Println(len(data[firstYamlDelim:]))
	if nextYamlDelim == -1 && atEOF == true {
		return 0, data[firstYamlDelim:], bufio.ErrFinalToken
	}
	if nextYamlDelim == -1 && atEOF == false {
		return 0, nil, nil
	}
	nextIndex := nextYamlDelim + index
	return nextIndex, data[firstYamlDelim:nextIndex], nil
}

func yamlSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	firstYamlDelim := bytes.Index(data, yamlDelim)
	index := firstYamlDelim + 3
	nextYamlDelim := bytes.Index(data[index:], yamlDelim)
	if nextYamlDelim == -1 && atEOF == true {
		return 0, data[firstYamlDelim:], bufio.ErrFinalToken
	}
	if nextYamlDelim == -1 && atEOF == false {
		return 0, nil, nil
	}
	nextIndex := nextYamlDelim + index
	return nextIndex, data[firstYamlDelim:nextIndex], nil
}

func main() {
	config := &config{
		StartBufferSize:  64 * 1024,
		MaxScanTokenSize: 64 * 1024 * 1024, //64MB max allocation size
	}

	multiconfig.MustLoad(config)

	if config.File == "" {
		log.Fatal("missing file to convert")
	}

	file, err := os.Open(config.File)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	buf := make([]byte, config.StartBufferSize)
	scanner.Buffer(buf, config.MaxScanTokenSize)

	scanner.Split(yamlSplit)

	jsonMap := map[string]interface{}{}

	for scanner.Scan() {
		yamlDocument := scanner.Text()

		//fmt.Println(yaml)
		names := name.FindStringSubmatch(yamlDocument)
		//TODO check error
		name := names[1]

		kinds := kind.FindStringSubmatch(yamlDocument)
		//TODO check error
		kind := kinds[1]
		fmt.Println(name, kind)

		if kind == "CustomResourceDefinition" || kind == "Namespace" {
			name = "0" + name
		}

		key := name + kind
		_, ok := jsonMap[key]
		if ok {
			log.Fatalf("dublicate found in jsonMap: %s", key)
		}

		var yamlObject interface{}
		err := yaml2.Unmarshal([]byte(yamlDocument), &yamlObject)
		if err != nil {
			log.Fatalf("failed to yaml2 unmarshal document: %s, raw block: %s", err, yamlDocument)
		}

		jsonMap[key] = dyno.ConvertMapI2MapS(yamlObject)
	}

	strings.Replace(file.Name(), ".yaml", ".json", 1)

	jsonFile, err := os.Create(fileName(file.Name()))
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}
	defer jsonFile.Close()

	err = json.NewEncoder(jsonFile).Encode(jsonMap)
	if err != nil {
		log.Fatalf("failed to encode jsonMap to json file: %s", err)
	}
	fmt.Println("json file has been written", jsonFile.Name())
}

func fileName(name string) string {
	if strings.HasSuffix(name, ".yaml") {
		return strings.Replace(name, ".yaml", ".json", 1)
	}
	if strings.HasSuffix(name, ".yml") {
		return strings.Replace(name, ".yml", ".json", 1)
	}
	return name + ".json"
}
