# yamltojson
Converts kubernetes yaml files to json files.

## installation

`go install github.com/xonvanetta/yamltojson@latest`

or this if you don't have go 1.17+

`go get github.com/xonvanetta/yamltojson`

## Usage
This will convert a kubernetes yaml file to one json file with keys based on `name+kind`

The command for converting files `./yamltojson -file tests/cert-manager.yaml`, the command will create a json file in `tests/cert-manager.json`

### Config
```
Usage of ./yamltojson:
  -file
    	Change value of File.
  -maxscantokensize
    	Change value of MaxScanTokenSize. (default 67108864)
  -startbuffersize
    	Change value of StartBufferSize. (default 65536)

Generated environment variables:
   CONFIG_FILE
   CONFIG_MAXSCANTOKENSIZE
   CONFIG_STARTBUFFERSIZE

flag: help requested
```


## Tests files
`cert-manager.yaml` is from https://github.com/jetstack/cert-manager/ release v1.5.4
