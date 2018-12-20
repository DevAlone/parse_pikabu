package parser

import (
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/op/go-logging"
	"github.com/streadway/amqp"
	"gogsweb.2-47.ru/d3dev/pikago"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

type NoTaskError struct{}

func (this NoTaskError) Error() string { return "there is no any task" }

type Parser struct {
	Config       *ParserConfig
	httpClient   *http.Client
	pikagoClient *pikago.Client
}

func NewParser() (*Parser, error) {
	parser := &Parser{}
	var err error
	parser.Config, err = NewParserConfigFromFile("")
	if err != nil {
		return nil, err
	}
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

	return parser, nil
}

func (this *Parser) handleError(err error) {
	if err == nil {
		panic("trying to handle nil error\n")
	} else if _, ok := err.(NoTaskError); ok {
		logger.ParserLog.Debug("there is no task, waiting...")
		time.Sleep(time.Duration(this.Config.WaitNoTaskSeconds) * time.Second)
		return
	}

	if e, ok := err.(*errors.Error); ok {
		logger.ParserLog.Error(e.ErrorStack())
	} else {
		logger.ParserLog.Error(err.Error())
	}

	time.Sleep(time.Duration(this.Config.WaitAfterErrorSeconds) * time.Second)
}

func (this *Parser) Loop() {
	for true {
		task, err := this.pullTask()
		if err != nil {
			this.handleError(err)
			continue
		}
		// process task
		err = this.processTask(task)
		if err != nil {
			this.handleError(err)
			continue
		}
	}
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

func (this *Parser) PutResultsToQueue(routingKey string, result interface{}) error {
	numberOfResults := 0
	resultType := reflect.TypeOf(result)
	switch resultType.Kind() {
	case reflect.Slice, reflect.Array:
		numberOfResults = reflect.ValueOf(result).Len()
	default:
		result = []interface{}{result}
		numberOfResults = 1
	}
	logger.ParserLog.Debugf("putting result to queue %v", result)

	var jsonMessage models.ParserResult
	jsonMessage.ParsingTimestamp = models.TimestampType(time.Now().Unix())
	jsonMessage.ParserId = "d3dev/" + this.Config.ParserId
	jsonMessage.NumberOfResults = numberOfResults
	jsonMessage.Results = result

	message, err := json.Marshal(jsonMessage)
	if err != nil {
		return err
	}

	conn, err := amqp.Dial(this.Config.AMQPAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"parser_results",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"parser_results",
		routingKey,
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         message,
		},
	)

	return err

}

func Main() {
	file, err := os.OpenFile("logs/parser.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	// loggingBackend := logger.NewLogBackend(os.Stderr, "", 0)
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, logger.LogFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)
	logger.ParserLog.Debug("app started")
	// TODO: pass config here

	parser, err := NewParser()
	if err != nil {
		panic(err)
	}
	parser.Loop()
}
