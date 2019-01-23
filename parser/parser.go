package parser

import (
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/parser/logger"
	"github.com/go-errors/errors"
	"github.com/op/go-logging"
	"gogsweb.2-47.ru/d3dev/pikago"
)

type NoTaskError struct{}

func (this NoTaskError) Error() string { return "there is no any task" }

type Parser struct {
	Config       *ParserConfig
	httpClient   *http.Client
	pikagoClient *pikago.Client
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
	)
	if err != nil {
		return nil, err
	}
	requestsSender, err := pikago.NewClientProxyRequestsSender(proxyProvider)
	requestsSender.SetTimeout(parser.Config.PikagoTimeout)
	if err != nil {
		return nil, err
	}
	parser.pikagoClient, err = pikago.NewClient(requestsSender)
	if err != nil {
		return nil, err
	}
	if config.Settings.Debug {
		pikago.LogsAreEnabled = true
		parser.pikagoClient.Log = logging.MustGetLogger("pikago")
		logging.SetLevel(logging.DEBUG, parser.pikagoClient.Log.Module)
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

/*
func (this *Parser) pullTask() (interface{}, error) {
	resp, err := this.doAPIRequest("get", "/get/tasks/any", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err)
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, NoTaskError{}
		case http.StatusUnauthorized:
			return nil, errors.Errorf("Unauthorized: %v", string(body))
		}

		return nil, errors.Errorf("%v", string(body))
	}

	var task struct {
		Name string           `json:"name"`
		Data *json.RawMessage `json:"data"`
	}
	err = pikago.JsonUnmarshal(body, &task)
	if err != nil {
		return nil, errors.New(err)
	}

	switch task.Name {
	case "parse_user_by_id":
		res := models.ParseUserByIdTask{}
		err = json.Unmarshal(*task.Data, &res)
		return res, err
	case "parse_user_by_username":
		res := models.ParseUserByUsernameTask{}
		err = json.Unmarshal(*task.Data, &res)
		return res, err
	case "simple":
		res := models.SimpleTask{}
		err = json.Unmarshal(*task.Data, &res)
		return res, err
	}

	return nil, errors.Errorf("bad task name: %v", task.Name)
}
*/

func Main() {
	file, err := os.OpenFile("logs/parser.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	panicOnError(err)
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, logger.LogFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)

	if config.Settings.Debug {
		logging.SetLevel(logging.DEBUG, "parse_pikabu/parser")
	} else {
		logging.SetLevel(logging.WARNING, "parse_pikabu/parser")
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
