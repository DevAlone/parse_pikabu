package parser

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-errors/errors"
)

// ParserConfig -
type ParserConfig struct {
	ParserID string
	// by default 1, can be used in config to
	// define multiple parsers with the same behavior
	// parser_id will be suffixed with number of copy
	NumberOfInstances    uint
	APIURL               string
	APITimeout           int
	APISessionID         string
	ProxyProviderAPIURL  string
	ProxyProviderTimeout int
	// everything else will be ignored if this value is not empty string
	FixedProxyAddress                string
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

// ParsersConfig -
type ParsersConfig struct {
	Configs []ParserConfig
}

// NewParserConfigFromBytes -
func NewParserConfigFromBytes(configData []byte) (*ParserConfig, error) {
	config := &ParserConfig{}

	config.ParserID = "unique_parser_id"
	config.NumberOfInstances = 1
	config.APIURL = "http://localhost:8080/api/v1"
	config.ProxyProviderAPIURL = ""
	config.ProxyProviderTimeout = 2 * 60
	config.PikagoTimeout = 60
	config.PikagoNumberOfRequestTries = 100
	config.PikagoChangeProxyOnNthBadTry = 5
	config.PikagoWaitBetweenProcessingPages = 1
	config.PikagoWaitBeforeNextRequestMs = 1000
	config.APITimeout = 60
	config.WaitAfterErrorSeconds = 1
	config.WaitNoTaskSeconds = 5 // TODO: delete?
	config.APISessionID = "put parser's session id here"
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
