package parser

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-errors/errors"
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
	PikagoNumberOfRequestTries       uint
	PikagoWaitBeforeNextRequestMs    uint
	PikagoChangeProxyOnNthBadTry     uint
	WaitAfterErrorSeconds            int
	WaitNoTaskSeconds                int
	AMQPAddress                      string
	LogHTTPQueries                   bool
	FileToStoreSSLKeys               string
	ProxyGettingPolicy               string
}

type ParsersConfig struct {
	Configs []ParserConfig
}

func NewParserConfigFromBytes(configData []byte) (*ParserConfig, error) {
	config := &ParserConfig{}

	config.ParserId = "unique_parser_id"
	config.NumberOfInstances = 1
	config.ApiURL = "http://localhost:8080/api/v1"
	config.ProxyProviderAPIURL = ""
	config.ProxyProviderTimeout = 2 * 60
	config.PikagoTimeout = 45
	config.PikagoNumberOfRequestTries = 25
	config.PikagoChangeProxyOnNthBadTry = 5
	config.PikagoWaitBetweenProcessingPages = 1
	config.PikagoWaitBeforeNextRequestMs = 1000
	config.ApiTimeout = 60
	config.WaitAfterErrorSeconds = 1
	config.WaitNoTaskSeconds = 5 // TODO: delete?
	config.ApiSessionId = "put parser's session id here"
	config.AMQPAddress = "amqp://guest:guest@localhost:5672"
	config.LogHTTPQueries = false
	config.FileToStoreSSLKeys = ""
	config.ProxyGettingPolicy = "ProxyGettingPoliceRandom"

	if len(configData) > 0 {
		err := json.Unmarshal([]byte(configData), config)
		if err != nil {
			return nil, errors.New(err)
		}
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
