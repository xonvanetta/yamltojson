package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/icza/dyno"
	yaml2 "sigs.k8s.io/yaml"
)

func main() {
	filePath := ""
	flag.StringVar(&filePath, "file", "", "yaml document to convert to json")
	flag.StringVar(&filePath, "f", "", "yaml document to convert to json")

	flag.Parse()
	if filePath == "" {
		log.Fatal("missing file to convert")
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	jsonMap := map[string]interface{}{}

	for _, yamlBytes := range bytes.Split(b, []byte("\n---")) {
		doc := &struct {
			Kind     string
			Metadata struct {
				Name string
			}
		}{}

		err := yaml2.Unmarshal(yamlBytes, doc)
		if err != nil {
			log.Fatalf("failed to unmrshal doc")
		}
		if doc.Kind == "" && doc.Metadata.Name == "" {
			continue
		}

		fmt.Println(doc.Metadata.Name, doc.Kind)
		name := doc.Metadata.Name
		if doc.Kind == "CustomResourceDefinition" || doc.Kind == "Namespace" {
			name = "0" + name
		}

		key := name + doc.Kind
		_, ok := jsonMap[key]
		if ok {
			log.Fatalf("dublicate found in jsonMap: %s", key)
		}

		var yamlObject interface{}
		err = yaml2.Unmarshal(yamlBytes, &yamlObject)
		if err != nil {
			log.Fatalf("failed to yaml2 unmarshal document: %s, raw block: %s", err, yamlBytes)
		}

		jsonMap[key] = dyno.ConvertMapI2MapS(yamlObject)
	}

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
