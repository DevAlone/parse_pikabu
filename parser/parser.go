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
	"strings"
	"time"
)

type ParserConfig struct {
	ApiURL string
}

type NoTaskError struct{}

func (this NoTaskError) Error() string { return "there is no any task" }

func NewParserConfigFromFile(filepath string) (*ParserConfig, error) {
	parserConfig := &ParserConfig{}
	parserConfig.ApiURL = "http://localhost:8080/api/v1"

	if len(filepath) > 0 {
		// open file
	}

	return parserConfig, nil
}

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
		Timeout: 10 * time.Second,
	}
	proxyProvider, err := pikago.GetProxyPyProxyProvider(
		"https://eivailohciihi4uquapach7abei9iesh.d3d.info/api/v1/",
		60,
	)
	if err != nil {
		return nil, err
	}
	requestsSender, err := pikago.NewClientProxyRequestsSender(proxyProvider)
	if err != nil {
		return nil, err
	}
	parser.pikagoClient, err = pikago.NewClient(requestsSender)
	if err != nil {
		return nil, err
	}

	return parser, nil
}

func handleError(err error) {
	if err == nil {
		panic("trying to handle nil error\n")
	} else if _, ok := err.(NoTaskError); ok {
		logger.ParserLog.Debug("there is no task, waiting...")
		time.Sleep(5 * time.Second)
		return
	}

	if e, ok := err.(*errors.Error); ok {
		logger.ParserLog.Error(e.ErrorStack())
	} else {
		logger.ParserLog.Error(err.Error())
	}

	time.Sleep(10 * time.Second)
}

func (this *Parser) Loop() {
	for true {
		task, err := this.pullTask()
		if err != nil {
			handleError(err)
			continue
		}
		// process task
		err = this.processTask(task)
		if err != nil {
			handleError(err)
			continue
		}
		time.Sleep(1 * time.Second)
	}
}

func (this *Parser) doAPIRequest(method string, url string, body io.Reader) (*http.Response, error) {
	method = strings.ToUpper(method)
	req, err := http.NewRequest(method, this.Config.ApiURL+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Session-Id", "Hohpuu8oogoShaituNoh8iebaesiYaeh")
	return this.httpClient.Do(req)
}

func (this *Parser) pullTask() (interface{}, error) {
	resp, err := this.doAPIRequest("get", "/get/tasks/any", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, NoTaskError{}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err)
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
	}

	return nil, errors.Errorf("bad task name: %v", task.Name)
}

func (this *Parser) PutResultToQueue(result interface{}) error {
	logger.ParserLog.Debugf("putting result to queue %v", result)

	var jsonMessage struct {
		ParsingTimestamp models.TimestampType `json:"parsing_timestamp"`
		ParserId         string               `json:"parser_id"`
		Data             interface{}          `json:"data"`
	}
	jsonMessage.ParsingTimestamp = models.TimestampType(time.Now().Unix())
	jsonMessage.ParserId = "d3dev/parser_id"
	jsonMessage.Data = result

	message, err := json.Marshal(jsonMessage)
	if err != nil {
		return err
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
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
		"user_profile",
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
