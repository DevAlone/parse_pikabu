package parser

import (
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"fmt"
	"github.com/go-errors/errors"
	"gogsweb.2-47.ru/d3dev/pikago"
	"io/ioutil"
	"net/http"
)

func (this *Parser) processTask(task interface{}) error {
	switch t := task.(type) {
	case models.ParseUserByIdTask:
		fmt.Printf("parse user by id task %v\n", t)
		resp, err := this.doAPIRequest(
			"get",
			"/take/parse_user_by_id_tasks/"+fmt.Sprint(t.PikabuId),
			nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		bytesResp, err := ioutil.ReadAll(resp.Body)
		textResp := string(bytesResp)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return errors.Errorf("unable to take task %v, error: %v", task, textResp)
		}
	case models.ParseUserByUsernameTask:
		fmt.Printf("parse user by username task %v\n", t)
		logger.Log.Debugf("taking task to parse user by username %v", t)
		resp, err := this.doAPIRequest(
			"get",
			"/take/parse_user_by_username_tasks/"+fmt.Sprint(t.Username),
			nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		bytesResp, err := ioutil.ReadAll(resp.Body)
		textResp := string(bytesResp)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return errors.Errorf("unable to take task %v, error: %s", task, textResp)
		}
		err = this.processParseUserByUsernameTask(t)
		if err != nil {
			return err
		}
	default:
		print("bad task %v\n", t)
	}
	return nil
}

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

	return this.PutResultToQueue(res)
}
