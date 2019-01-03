package parser

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"io/ioutil"
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

func NewParserConfigFromBytes(configData []byte) (*ParserConfig, error) {
	config := &ParserConfig{}

	config.ParserId = "unique_parser_id"
	config.NumberOfInstances = 1
	config.ApiURL = "http://localhost:8080/api/v1"
	config.ProxyProviderAPIURL = "https://eivailohciihi4uquapach7abei9iesh.d3d.info/api/v2/"
	config.ProxyProviderTimeout = 60
	config.PikagoTimeout = 30
	config.PikagoWaitBetweenProcessingPages = 1
	config.ApiTimeout = 60
	config.WaitAfterErrorSeconds = 10
	config.WaitNoTaskSeconds = 5
	config.ApiSessionId = "parser_oogoShaituNoh8iebaesiYaeh"
	config.AMQPAddress = "amqp://guest:guest@localhost:5672"

	err := json.Unmarshal([]byte(configData), config)
	if err != nil {
		return nil, errors.New(err)
	}

	return config, nil
}

func NewParsersConfigFromFile(filename string) (*ParsersConfig, error) {
	// parsersConfig := &ParsersConfig{}
	var jsonData struct {
		Configs []json.RawMessage
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &jsonData)

	parsersConfig := &ParsersConfig{}
	for _, item := range jsonData.Configs {
		bytes, err := item.MarshalJSON()
		if err != nil {
			return nil, err
		}
		config, err := NewParserConfigFromBytes(bytes)
		if err != nil {
			return nil, err
		}
		parsersConfig.Configs = append(parsersConfig.Configs, *config)
	}

	return parsersConfig, err
}