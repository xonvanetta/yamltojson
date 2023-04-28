# yamltojson
Converts kubernetes yaml files to json files.

## installation

`go install github.com/xonvanetta/yamltojson@latest`

or this if you don't have go 1.17+

`go get github.com/xonvanetta/yamltojson`

## Usage
This will convert a kubernetes yaml file to one json file with keys based on `name+kind`

The command for converting files `yamltojson -file tests/cert-manager.yaml`, the command will create a json file in `tests/cert-manager.json`

### Config
```
$ yamltojson --help
Usage of yamltojson:
  -file value
    	Change value of File.
  -namepattern value
    	Change value of NamePattern. (default {{if or (eq .Kind "CustomResourceDefinition") (eq .Kind "Namespace")}}0{{end}}{{ .Metadata.Name }}-{{ .Kind }})
  -seperatecrds
    	Change value of SeperateCRDs. (default false)

Generated environment variables:
   CONFIG_FILE
   CONFIG_NAMEPATTERN
   CONFIG_SEPERATECRDS
```


## Tests files
`cert-manager.yaml` is from https://github.com/jetstack/cert-manager/ release v1.5.4
`fluxv2.yaml` is from https://github.com/fluxcd/flux2/releases/download/v2.0.0-rc.1/install.yaml release v2.0.0-rc.1
