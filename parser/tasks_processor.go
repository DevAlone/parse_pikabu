package parser

import (
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"gogsweb.2-47.ru/d3dev/pikago"
	"time"
)

func (this *Parser) processTask(task interface{}) error {
	switch t := task.(type) {
	case models.ParseUserByIdTask:
		//logger.Log.Debugf("taking task to parse user by id %v", t)
		//if err := this.takeTask("parse_user_by_id_tasks", t.Id); err != nil {
		//	return err
		//}
		// TODO: process
	case models.ParseUserByUsernameTask:
		//logger.Log.Debugf("taking task to parse user by username %v", t)
		//if err := this.takeTask("parse_user_by_username_tasks", t.Id); err != nil {
		//	return err
		//}
		err := this.processParseUserByUsernameTask(t)
		if err != nil {
			return err
		}
	case models.SimpleTask:
		//logger.Log.Debugf("taking simple task %v", t)
		//if err := this.takeTask("simple_tasks", t.Id); err != nil {
		//	return err
		//}
		return this.processSimpleTask(t)
	default:
		return errors.Errorf("bad task %v\n", t)
	}
	return nil
}

//func (this *Parser) takeTask(
//	taskTableName string,
//	idField uint64,
//) error {
//	resp, err := this.doAPIRequest(
//		"get",
//		"/take/"+taskTableName+"/"+fmt.Sprint(idField),
//		nil)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//	bytesResp, err := ioutil.ReadAll(resp.Body)
//	textResp := string(bytesResp)
//	if err != nil {
//		return err
//	}
//	if resp.StatusCode != http.StatusOK {
//		return errors.Errorf("unable to take task %v, error: %s", taskTableName, textResp)
//	}
//
//	return nil
//}

func (this *Parser) processParseUserByUsernameTask(task models.ParseUserByUsernameTask) error {
	logger.Log.Debugf("sending request to get user %v", task.Username)
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
		logger.Log.Debugf("sending request to get communities")
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
