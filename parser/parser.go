package parser

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/parser/logger"
	"github.com/go-errors/errors"
	"gogsweb.2-47.ru/d3dev/pikago"
)

// NoTaskError -
type NoTaskError struct{}

func (e NoTaskError) Error() string { return "there is no any task" }

// Parser -
type Parser struct {
	Config       *ParserConfig
	httpClient   *http.Client
	pikagoClient *pikago.MobileClient
}

// NewParser - creates new parser
func NewParser(parserConfig *ParserConfig) (*Parser, error) {
	parser := &Parser{}
	var err error
	parser.Config = parserConfig
	parser.httpClient = &http.Client{
		Timeout: time.Duration(parser.Config.APITimeout) * time.Second,
	}

	proxyGettingPolicy := pikago.ProxyGettingPoliceRandom
	switch parserConfig.ProxyGettingPolicy {
	case "ProxyGettingPoliceRandom":
		proxyGettingPolicy = pikago.ProxyGettingPoliceRandom
	case "ProxyGettingPolice1024BestResponseTime":
		proxyGettingPolicy = pikago.ProxyGettingPolice1024BestResponseTime
	default:
		return nil, errors.Errorf("bad proxy getting policy %v", parserConfig.ProxyGettingPolicy)
	}

	proxyProvider, err := pikago.GetProxyPyProxyProvider(
		parser.Config.ProxyProviderAPIURL,
		parser.Config.ProxyProviderTimeout,
		proxyGettingPolicy,
	)

	if err != nil {
		return nil, err
	}
	requestsSender, err := pikago.NewClientProxyRequestsSender(proxyProvider)
	if err != nil {
		return nil, err
	}
	requestsSender.NumberOfRequestTries = parser.Config.PikagoNumberOfRequestTries
	requestsSender.ChangeProxyOnNthBadTry = parser.Config.PikagoChangeProxyOnNthBadTry
	requestsSender.WaitBeforeNextRequestMs = parser.Config.PikagoWaitBeforeNextRequestMs
	requestsSender.SetTimeout(time.Duration(parser.Config.PikagoTimeout) * time.Second)
	if len(parser.Config.FileToStoreSSLKeys) > 0 {
		err := requestsSender.SetFileToStoreSSLKeys(parser.Config.FileToStoreSSLKeys)
		if err != nil {
			return nil, err
		}
	}
	parser.pikagoClient, err = pikago.NewClient(requestsSender)
	if err != nil {
		return nil, err
	}
	parser.pikagoClient.SetLog(logger.PikagoLog, logger.PikagoHttpLog)
	parser.pikagoClient.AddBeforeRequestMiddleware(func(req *http.Request) *http.Request {
		_ = parser.pikagoClient.ResetState()
		return req
	})

	return parser, nil
}

func (p *Parser) handleError(err error) {
	if err == nil {
		return
	}

	if _, ok := err.(NoTaskError); ok {
		// logger.ParserLog.Debug("there is no task, waiting...")
		time.Sleep(time.Duration(p.Config.WaitNoTaskSeconds) * time.Second)
		return
	}

	logger.Log.Error(helpers.ErrorToString(err))

	time.Sleep(time.Duration(p.Config.WaitAfterErrorSeconds) * time.Second)
}

// TODO: consider removing
func (p *Parser) doAPIRequest(method string, url string, body io.Reader) (*http.Response, error) {
	method = strings.ToUpper(method)
	req, err := http.NewRequest(method, p.Config.APIURL+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Session-Id", p.Config.APISessionID)
	return p.httpClient.Do(req)
}
