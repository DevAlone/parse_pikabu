package parser

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"
	"gogsweb.2-47.ru/d3dev/pikago"
)

func (p *Parser) processParseUserTask(task *models.ParseUserTask) error {
	var res *models.ParserUserProfileResultData
	var err error

	if len(task.Username) > 0 {
		res, err = p.processParseUserTaskByUsername(task)
		if err != nil {
			if _, ok := err.(*pikago.PikabuErrorRequestedPageNotFound); ok {
				res, err = p.ProcessParseUserTaskByID(task)
			} else {
				return errors.Errorf("Error while processing task %v. Error: %v", task, err)
			}
		}
	} else {
		res, err = p.ProcessParseUserTaskByID(task)
	}

	if err != nil {
		if pe, ok := err.(*pikago.PikabuErrorRequestedPageNotFound); ok {
			return p.PutResultsToQueue("user_profile_not_found", []models.ParserUserProfileNotFoundResultData{
				{
					PikabuID:    task.PikabuID,
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

// ProcessParseUserTaskByID -
func (p *Parser) ProcessParseUserTaskByID(task *models.ParseUserTask) (*models.ParserUserProfileResultData, error) {
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

	httpReq, err := makeRequest(task.PikabuID)
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
		}
		return nil, err
	} else if err != nil {
		return nil, err
	}

	var resp struct {
		Result      bool   `json:"result"`
		Message     string `json:"message"`
		MessageCode int    `json:"message_code"`
		Data        struct {
			HTML string `json:"html"`
		} `json:"data"`
	}
	err = pikago.JsonUnmarshal(body, &resp, true)
	if err != nil {
		return nil, err
	}
	if !resp.Result {
		return nil, pikago.NewPikabuError(fmt.Sprintf("req: %v, resp: %v", httpReq, resp))
	}

	resp.Data.HTML = strings.Replace(resp.Data.HTML, "\n", "", -1)

	regex, err := regexp.Compile(`@.+?"`)
	if err != nil {
		return nil, err
	}
	username := regex.FindString(resp.Data.HTML)
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
