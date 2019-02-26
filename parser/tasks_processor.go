package parser

import (
	go_errors "errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/globals"

	"gogsweb.2-47.ru/d3dev/pikago"

	"bitbucket.org/d3dev/parse_pikabu/amqp_helper"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/parser/logger"
	"github.com/go-errors/errors"
	"github.com/streadway/amqp"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

func (p *Parser) Loop() {
	for true {
		err := p.ListenForTasks()
		if err != nil {
			p.handleError(err)
			time.Sleep(time.Duration(p.Config.WaitAfterErrorSeconds) * time.Second)
		}
	}
}

func (p *Parser) ListenForTasks() error {
	defer func() {
		if r := recover(); r != nil {
			p.handleError(errors.Errorf("panic: %v", r))
		}
	}()

	if globals.SingleProcessMode {
		return p.ListenForChannelTasks()
	} else {
		return p.ListenForAMQPTasks()
	}
}

func (p *Parser) ListenForChannelTasks() error {
	for true {
		select {
		case task := <-globals.ParserParseUserTasks:
			return p.processParseUserTask(task)
		case task := <-globals.ParserSimpleTasks:
			switch task.Name {
			case "parse_communities_pages":
				return p.processParseCommunitiesPagesTask()
			default:
				return errors.Errorf("unknown type of simple task %v", task)
			}
		}
	}

	return nil
}

func (p *Parser) ListenForAMQPTasks() error {
	logger.Log.Debug("connecting to amqp server...")
	connection, err := amqp_helper.GetAMQPConnection(p.Config.AMQPAddress)
	if err != nil {
		return errors.New(err)
	}

	ch, err := connection.Channel()
	if err != nil {
		return errors.New(err)
	}
	defer func() { _ = ch.Close() }()

	err = amqp_helper.DeclareExchanges(ch)
	if err != nil {
		return errors.New(err)
	}

	q, err := amqp_helper.DeclareParserTasksQueue(ch)
	if err != nil {
		return errors.New(err)
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"parser_tasks",
		false,
		nil,
	)
	if err != nil {
		return errors.New(err)
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return errors.New(err)
	}

	messages, err := ch.Consume(
		q.Name,
		// TODO: process only those tasks that can be processed by p parser
		"", // routing key
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New(err)
	}

	logger.Log.Debug("start waiting for tasks")
	defer logger.Log.Debug("stop waiting for parser results")

	for message := range messages {
		err = message.Ack(false)
		if err != nil {
			return err
		}
		err = p.processMessage(message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) processMessage(message amqp.Delivery) error {
	logger.Log.Debugf("got message: %v", string(message.Body))

	switch message.RoutingKey {
	case "parse_user":
		var task models.ParseUserTask
		err := pikago.JsonUnmarshal(message.Body, &task, true)
		if err != nil {
			return errors.New(err)
		}
		return p.processParseUserTask(&task)
	case "parse_communities_pages":
		return p.processParseCommunitiesPagesTask()
	default:
		logger.Log.Warningf(
			"Unregistered task type \"%v\". Message: \"%v\"",
			message.RoutingKey,
			string(message.Body),
		)
		return nil
	}
}

func (p *Parser) processParseUserTask(task *models.ParseUserTask) error {
	var res *models.ParserUserProfileResultData
	var err error

	if len(task.Username) > 0 {
		res, err = p.processParseUserTaskByUsername(task)
		if _, ok := err.(*pikago.PikabuErrorRequestedPageNotFound); err != nil && ok {
			res, err = p.ProcessParseUserTaskById(task)
		} else if err != nil {
			return go_errors.New(fmt.Sprintf("Error while processing task %v. Error: %v", task, err))
		}
	} else {
		res, err = p.ProcessParseUserTaskById(task)
	}

	if err != nil {
		if pe, ok := err.(*pikago.PikabuErrorRequestedPageNotFound); ok {
			return p.PutResultsToQueue("user_profile_not_found", []models.ParserUserProfileNotFoundResultData{
				{
					PikabuId:    task.PikabuId,
					Username:    task.Username,
					PikabuError: pe,
				},
			})
		}
		return err
	}

	return p.PutResultsToQueue("user_profile", []models.ParserUserProfileResultData{
		*res,
	})
}

func (p *Parser) ProcessParseUserTaskById(task *models.ParseUserTask) (*models.ParserUserProfileResultData, error) {
	// parse by id
	makeRequest := func(id uint64) (*http.Request, error) {
		url := fmt.Sprintf("https://pikabu.ru/ajax/user_info.php?action=get_short_profile&user_id=%v", id)
		httpReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		token := helpers.GetRandomString([]rune("abcdefghijklmnopqrstuvwxyz0123456789"), 32)
		httpReq.Header.Add("X-Csrf-Token", token)
		httpReq.Header.Add("Cookie", fmt.Sprintf("PHPSESS=%v;", token))
		httpReq.Header.Add("X-Requested-With", "XMLHttpRequest")

		return httpReq, nil
	}

	httpReq, err := makeRequest(task.PikabuId)
	if err != nil {
		return nil, err
	}

	body, _, err := p.pikagoClient.DoHttpRequest(httpReq)
	if pikabuErr, ok := err.(*pikago.PikabuError); ok && pikabuErr.StatusCode == 403 {
		// check that 403 was because user does not exist
		httpReq, e := makeRequest(1)
		if e != nil {
			return nil, e
		}
		_, resp, e := p.pikagoClient.DoHttpRequest(httpReq)
		// TODO: check actual data
		if resp.StatusCode == 200 {
			return nil, &pikago.PikabuErrorRequestedPageNotFound{
				PikabuError: pikago.PikabuError{
					StatusCode: 404,
					Message:    "Not found by id",
				},
			}
		} else {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	var resp struct {
		Result      bool   `json:"result"`
		Message     string `json:"message"`
		MessageCode int    `json:"message_code"`
		Data        struct {
			Html string `json:"html"`
		} `json:"data"`
	}
	err = pikago.JsonUnmarshal(body, &resp, true)
	if err != nil {
		return nil, err
	}
	if resp.Result == false {
		return nil, pikago.NewPikabuError(fmt.Sprintf("req: %v, resp: %v", httpReq, resp))
	}

	resp.Data.Html = strings.Replace(resp.Data.Html, "\n", "", -1)

	regex, err := regexp.Compile(`@.+?"`)
	if err != nil {
		return nil, err
	}
	username := regex.FindString(resp.Data.Html)
	if len(username) < 3 {
		return nil, errors.Errorf("Unable to find username in data. Resp: %v", resp)
	}
	username = username[1 : len(username)-1]

	task.Username = username
	return p.processParseUserTaskByUsername(task)
}

func (p *Parser) processParseUserTaskByUsername(task *models.ParseUserTask) (*models.ParserUserProfileResultData, error) {
	userProfile, err := p.pikagoClient.UserProfileGet(task.Username)
	if err != nil {
		return nil, err
	}
	res := models.ParserUserProfileResultData{}
	res.User = userProfile

	return &res, nil
}

func (p *Parser) processParseCommunitiesPagesTask() error {
	results := []pikago_models.CommunitiesPage{}

	page := 0
	for true {
		communitiesPage, err := p.pikagoClient.CommunitiesGet(page)
		if err != nil {
			return err
		}
		if len(communitiesPage.List) == 0 {
			break
		}
		results = append(results, *communitiesPage)

		page++
		time.Sleep(time.Duration(p.Config.PikagoWaitBetweenProcessingPages) * time.Second)
	}

	return p.PutResultsToQueue("communities_pages", results)
}
