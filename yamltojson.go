package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"

	"github.com/icza/dyno"
	"github.com/koding/multiconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	NamePattern  string
	SeperateCRDs bool
	File         string
}

func main() {
	config := &Config{
		NamePattern: `{{if or (eq .Kind "CustomResourceDefinition") (eq .Kind "Namespace")}}0{{end}}{{ .Metadata.Name }}-{{ .Kind }}`,
	}

	multiconfig.MustLoad(config)

	filePath := config.File
	if filePath == "" {
		log.Fatal("missing first arg for file")
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

	tmpl, err := template.New("name").Parse(config.NamePattern)
	if err != nil {
		log.Fatalf("failed to parse template for name: %s", err)
	}
	jsonMap := map[string]interface{}{}
	crdMap := map[string]interface{}{}

	for _, yamlBytes := range bytes.Split(b, []byte("\n---")) {
		doc := &struct { //maybe map string to easier use with pattern
			Kind     string
			Metadata struct {
				Name      string
				Namespace string
			}
		}{}

		err := yaml.Unmarshal(yamlBytes, doc)
		if err != nil {
			log.Fatalf("failed to unmrshal doc")
		}
		if doc.Kind == "" && doc.Metadata.Name == "" {
			continue
		}

		fmt.Println(doc.Metadata.Name, doc.Kind)

		b := &bytes.Buffer{}
		err = tmpl.Execute(b, doc)
		if err != nil {
			log.Fatalf("failed to execute template: %s", err)
		}

		key := b.String()
		_, ok := jsonMap[key]
		if ok {
			log.Fatalf("dublicate found in jsonMap: %s", key)
		}

		var yamlObject interface{}
		err = yaml.Unmarshal(yamlBytes, &yamlObject)
		if err != nil {
			log.Fatalf("failed to yaml2 unmarshal document: %s, raw block: %s", err, yamlBytes)
		}

		if config.SeperateCRDs && doc.Kind == "CustomResourceDefinition" {
			crdMap[key] = dyno.ConvertMapI2MapS(yamlObject)
		} else {
			jsonMap[key] = dyno.ConvertMapI2MapS(yamlObject)
		}
	}
	err = saveAsJson(replaceYamlToJson(file.Name())+".json", jsonMap)
	if err != nil {
		log.Fatalf("failed to save as json: %s", err)
	}

	if config.SeperateCRDs {
		err := saveAsJson(replaceYamlToJson(file.Name()+"-crd.json"), crdMap)
		if err != nil {
			log.Fatalf("failed to save as json: %s", err)
		}
	}
}

func saveAsJson(name string, data interface{}) error {
	jsonFile, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}
	defer jsonFile.Close()

	err = json.NewEncoder(jsonFile).Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode jsonMap to json file: %s", err)
	}
	fmt.Println("json file has been written", jsonFile.Name())
	return nil
}

func replaceYamlToJson(name string) string {
	name = strings.Replace(name, ".yaml", "", 1)
	name = strings.Replace(name, ".yml", "", 1)
	return name
}
