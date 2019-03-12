package parser

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/streadway/amqp"

	"bitbucket.org/d3dev/parse_pikabu/parser/logger"
	"github.com/go-errors/errors"
	"gogsweb.2-47.ru/d3dev/pikago"
)

type NoTaskError struct{}

func (this NoTaskError) Error() string { return "there is no any task" }

type Parser struct {
	Config       *ParserConfig
	httpClient   *http.Client
	pikagoClient *pikago.MobileClient
	amqpChannel  *amqp.Channel
}

func NewParser(parserConfig *ParserConfig) (*Parser, error) {
	parser := &Parser{}
	var err error
	parser.Config = parserConfig
	parser.httpClient = &http.Client{
		Timeout: time.Duration(parser.Config.ApiTimeout) * time.Second,
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

func (this *Parser) handleError(err error) {
	if err == nil {
		return
	}

	if _, ok := err.(NoTaskError); ok {
		// logger.ParserLog.Debug("there is no task, waiting...")
		time.Sleep(time.Duration(this.Config.WaitNoTaskSeconds) * time.Second)
		return
	}

	if e, ok := err.(*errors.Error); ok {
		logger.Log.Error(e.ErrorStack())
	} else {
		logger.Log.Error(err.Error())
	}

	time.Sleep(time.Duration(this.Config.WaitAfterErrorSeconds) * time.Second)
}

func (this *Parser) doAPIRequest(method string, url string, body io.Reader) (*http.Response, error) {
	method = strings.ToUpper(method)
	req, err := http.NewRequest(method, this.Config.ApiURL+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Session-Id", this.Config.ApiSessionId)
	return this.httpClient.Do(req)
}
