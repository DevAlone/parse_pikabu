package parser

import (
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"gogsweb.2-47.ru/d3dev/pikago"
	"time"
)

func (this *Parser) processTask(task interface{}) error {
	switch t := task.(type) {
	case models.ParseUserByIdTask:
		err := this.processParseUserByIdTask(t)
		if err != nil {
			return err
		}
	case models.ParseUserByUsernameTask:
		err := this.processParseUserByUsernameTask(t)
		if err != nil {
			return err
		}
	case models.SimpleTask:
		return this.processSimpleTask(t)
	default:
		return errors.Errorf("bad task %v\n", t)
	}
	return nil
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

func (this *Parser) processSimpleTask(task models.SimpleTask) error {
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
	return this.PutResultsToQueue("communities_page", results)
}
