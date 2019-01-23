package parser

import (
	"bitbucket.org/d3dev/parse_pikabu/amqp_helper"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/parser/logger"
	go_errors "errors"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/streadway/amqp"
	"gogsweb.2-47.ru/d3dev/pikago"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func (this *Parser) Loop() {
	for true {
		err := this.ListenForTasks()
		if err != nil {
			this.handleError(err)
			time.Sleep(time.Duration(this.Config.WaitAfterErrorSeconds) * time.Second)
		}
	}
}

func (this *Parser) ListenForTasks() error {
	defer func() {
		if r := recover(); r != nil {
			this.handleError(errors.Errorf("panic: %v", r))
		}
	}()

	logger.Log.Debug("connecting to amqp server...")
	connection, err := amqp_helper.GetAMQPConnection(this.Config.AMQPAddress)
	if err != nil {
		return errors.New(err)
	}

	ch, err := connection.Channel()
	if err != nil {
		return errors.New(err)
	}
	defer func() { _ = ch.Close() }()

	err = ch.ExchangeDeclare(
		"parser_tasks",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New(err)
	}

	// TODO: move to another file
	q, err := ch.QueueDeclare(
		"bitbucket.org/d3dev/parse_pikabu/parser_tasks",
		true,
		false,
		false,
		false,
		nil,
		/*
			amqp.Table{
				"x-queue-mode": "lazy",
			},
		*/
	)
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
		// TODO: process only those tasks that can be processed by this parser
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
		err = this.processMessage(message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Parser) processMessage(message amqp.Delivery) error {
	logger.Log.Debugf("got message: %v", string(message.Body))

	switch message.RoutingKey {
	case "parse_user":
		var task models.ParseUserTask
		err := pikago.JsonUnmarshal(message.Body, &task)
		if err != nil {
			return errors.New(err)
		}
		return this.processParseUserTask(task)
	case "parse_communities_pages":
		return this.processParseCommunitiesPagesTask()
	default:
		logger.Log.Warningf(
			"Unregistered task type \"%v\". Message: \"%v\"",
			message.RoutingKey,
			string(message.Body),
		)
		return nil
	}
}

func (this *Parser) processParseUserTask(task models.ParseUserTask) error {
	var res *struct {
		User *pikago.UserProfile `json:"user"`
	}
	var err error

	if len(task.Username) > 0 {
		res, err = this.processParseUserTaskByUsername(task)
		if pikabuErr, ok := err.(pikago.PikabuError); err != nil && ok && strings.Contains(pikabuErr.Message, "could not be found") {
			res, err = this.processParseUserTaskById(task)
			if err != nil {
				return err
			}
		} else if err != nil {
			return go_errors.New(fmt.Sprintf("Error while processing task %v. Error: %v", task, err))
		}
	} else {
		res, err = this.processParseUserTaskById(task)
		if err != nil {
			return err
		}
	}

	return this.PutResultsToQueue("user_profile", res)
}

func (this *Parser) processParseUserTaskById(task models.ParseUserTask) (*struct {
	User *pikago.UserProfile `json:"user"`
}, error) {
	// parse by id
	url := fmt.Sprintf("https://pikabu.ru/ajax/user_info.php?action=get_short_profile&user_id=%v", task.PikabuId)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	token := helpers.GetRandomString([]rune("abcdefghijklmnopqrstuvwxyz0123456789"), 32)
	httpReq.Header.Add("X-Csrf-Token", token)
	httpReq.Header.Add("Cookie", fmt.Sprintf("PHPSESS=%v;", token))
	httpReq.Header.Add("X-Requested-With", "XMLHttpRequest")

	body, httpResp, err := this.pikagoClient.DoHttpRequest(httpReq)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != 403 {
		// TODO: process deleted users somehow
	}

	var resp struct {
		Result      bool   `json:"result"`
		Message     string `json:"message"`
		MessageCode int    `json:"message_code"`
		Data        struct {
			Html string `json:"html"`
		} `json:"data"`
	}
	err = pikago.JsonUnmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Result == false {
		return nil, pikago.NewPikabuError(fmt.Sprintf("req: %v, resp: %v", httpReq, resp))
	}

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
	return this.processParseUserTaskByUsername(task)
}

func (this *Parser) processParseUserTaskByUsername(task models.ParseUserTask) (*struct {
	User *pikago.UserProfile `json:"user"`
}, error) {
	userProfile, err := this.pikagoClient.UserProfileGet(task.Username)
	if err != nil {
		return nil, err
	}
	var res struct {
		User *pikago.UserProfile `json:"user"`
	}
	res.User = userProfile

	return &res, nil
}

func (this *Parser) processParseCommunitiesPagesTask() error {
	results := []pikago.CommunitiesPage{}

	page := 0
	for true {
		communitiesPage, err := this.pikagoClient.CommunitiesGet(page)
		if err != nil {
			return err
		}
		if len(communitiesPage.List) == 0 {
			break
		}
		results = append(results, *communitiesPage)

		page++
		time.Sleep(time.Duration(this.Config.PikagoWaitBetweenProcessingPages) * time.Second)
	}

	return this.PutResultsToQueue("communities_pages", results)
}
