package parser

import (
	"encoding/json"
	"os"
)

type ParserConfig struct {
	ParserId string
	// by default 1, can be used in config to
	// define multiple parsers with the same behavior
	// parser_id will be suffixed with number of copy
	NumberOfInstances                uint
	ApiURL                           string
	ApiTimeout                       int
	ApiSessionId                     string
	ProxyProviderAPIURL              string
	ProxyProviderTimeout             int
	PikagoTimeout                    uint
	PikagoWaitBetweenProcessingPages int
	WaitAfterErrorSeconds            int
	WaitNoTaskSeconds                int
	AMQPAddress                      string
}

type ParsersConfig struct {
	Configs []ParserConfig
}

func NewParserConfigFromString(configData string) (*ParserConfig, error) {
	config := &ParserConfig{}

	config.ParserId = "unique_parser_id"
	config.NumberOfInstances = 1
	config.ApiURL = "http://localhost:8080/api/v1"
	config.ProxyProviderAPIURL = "https://eivailohciihi4uquapach7abei9iesh.d3d.info/api/v1/"
	config.ProxyProviderTimeout = 60
	config.PikagoTimeout = 5
	config.PikagoWaitBetweenProcessingPages = 1
	config.ApiTimeout = 60
	config.WaitAfterErrorSeconds = 5
	config.WaitNoTaskSeconds = 2
	config.ApiSessionId = "parser_oogoShaituNoh8iebaesiYaeh"
	config.AMQPAddress = "amqp://guest:guest@localhost:5672"

	err := json.Unmarshal([]byte(configData), config)

	return config, err
}

func NewParsersConfigFromFile(filename string) (*ParsersConfig, error) {
	parsersConfig := &ParsersConfig{}

	file, err := os.Open("parsers.config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(parsersConfig)

	return parsersConfig, nil
}
