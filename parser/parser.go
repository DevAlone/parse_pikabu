package parser

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"github.com/streadway/amqp"

	"bitbucket.org/d3dev/parse_pikabu/parser/logger"
	"github.com/go-errors/errors"
	logging "github.com/op/go-logging"
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
	proxyProvider, err := pikago.GetProxyPyProxyProvider(
		parser.Config.ProxyProviderAPIURL,
		parser.Config.ProxyProviderTimeout,
		pikago.ProxyGettingPoliceRandom,
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
	parser.pikagoClient, err = pikago.NewClient(requestsSender)
	if err != nil {
		return nil, err
	}
	parser.pikagoClient.AddBeforeRequestMiddleware(func(req *http.Request) *http.Request {
		_ = parser.pikagoClient.ResetState()
		return req
	})
	if config.Settings.Debug {
		logging.SetLevel(logging.DEBUG, "pikago")
	} else {
		logging.SetLevel(logging.WARNING, "pikago")
	}

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

func Main() {
	file, err := os.OpenFile("logs/parser.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	panicOnError(err)
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, logger.LogFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)

	if config.Settings.Debug {
		logging.SetLevel(logging.DEBUG, "parse_pikabu")
	} else {
		logging.SetLevel(logging.WARNING, "parse_pikabu")
	}

	logger.Log.Debug("parsers started")

	parsersConfig, err := NewParsersConfigFromFile("parsers.config.json")
	panicOnError(err)

	var wg sync.WaitGroup

	for _, parserConfig := range parsersConfig.Configs {
		// var configs
		for i := uint(0); i < parserConfig.NumberOfInstances; i++ {
			var conf ParserConfig
			conf = parserConfig
			if i != 0 {
				conf.ParserId += "_copy_" + fmt.Sprint(i)
			}

			parser, err := NewParser(&conf)
			if err != nil {
				panicOnError(err)
			}
			wg.Add(1)
			go func() {
				parser.Loop()
				wg.Done()
				err := parser.Cleanup()
				panicOnError(err)
			}()
		}
	}

	wg.Wait()
}

func panicOnError(err error) {
	if err == nil {
		return
	}
	if e, ok := err.(*errors.Error); ok {
		_, _ = os.Stderr.WriteString(e.ErrorStack())
	}

	panic(err)
}
