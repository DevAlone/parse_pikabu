package parser

import (
	"bitbucket.org/d3dev/parse_pikabu/amqp_helper"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/parser/logger"
	"github.com/go-errors/errors"
	"github.com/streadway/amqp"
	"gogsweb.2-47.ru/d3dev/pikago"
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
		err = this.processMessage(message)
		if err != nil {
			return err
		}
		err = message.Ack(false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Parser) processMessage(message amqp.Delivery) error {
	logger.Log.Debugf("got message: %v", string(message.Body))

	switch message.RoutingKey {
	case "parse_user_by_username":
		var task models.ParseUserByUsernameTask
		err := pikago.JsonUnmarshal(message.Body, &task)
		if err != nil {
			return errors.New(err)
		}
		return this.processParseUserByUsernameTask(task)
	case "parse_user_by_id":
		var task models.ParseUserByIdTask
		err := pikago.JsonUnmarshal(message.Body, &task)
		if err != nil {
			return errors.New(err)
		}
		return this.processParseUserByIdTask(task)
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

func (this *Parser) processParseUserByIdTask(task models.ParseUserByIdTask) error {
	/*
		curl -v 'https://pikabu.ru/ajax/user_info.php?action=get_short_profile&user_id=1'
			-H 'X-Csrf-Token: 89hvsja20e8ivco081oboj6fgnfpmq45'
			-H 'X-Requested-With: XMLHttpRequest'
			-H 'Cookie: PHPSESS=89hvsja20e8ivco081oboj6fgnfpmq45;'
	*/
	// TODO: complete
	return nil
}

func (this *Parser) processParseUserByUsernameTask(task models.ParseUserByUsernameTask) error {
	userProfile, err := this.pikagoClient.UserProfileGet(task.Username)
	if err != nil {
		return err
	}
	var res struct {
		User *pikago.UserProfile `json:"user"`
	}
	res.User = userProfile

	return this.PutResultsToQueue("user_profile", res)
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
