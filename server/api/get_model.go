package api

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"gogsweb.2-47.ru/d3dev/pikago"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var availableModels = map[string]interface{}{
	"pikabu_user": []models.PikabuUser{},
}

func GetModel(c *gin.Context) {
	var request struct {
		Name          string        `json:"name"`
		OrderByFields string        `json:"order_by_fields"`
		Offset        pikago.UInt64 `json:"offset"`
		Limit         pikago.UInt64 `json:"limit"`
	}

	err := c.Bind(&request)
	if err != nil {
		fmt.Println(err)
		logger.Log.Debug("error: ", err)
		AnswerError(c, http.StatusBadRequest, "your request is sooo bad")
		return
	}

	if request.Limit.Value == 0 {
		request.Limit.Value = uint64(config.Settings.ServerMaximumNumberOfResultsPerPage)
	}

	if request.Limit.Value > uint64(config.Settings.ServerMaximumNumberOfResultsPerPage) {
		AnswerError(c, http.StatusBadRequest, "you want too many of it")
		return
	}

	model, found := availableModels[request.Name]
	if !found {
		AnswerError(
			c,
			http.StatusBadRequest,
			"there is not any model like this, or you're not allowed to see it",
		)
		return
	}

	typeOfResult := reflect.TypeOf(model).Elem()
	results := reflect.New(reflect.TypeOf(model)).Interface()

	dbReq := models.Db.Model(results)
	orderBy, err := orderByFieldsToGoPg(request.OrderByFields, typeOfResult)
	if err != nil {
		logger.Log.Debug("error: ", err)
		AnswerError(c, http.StatusBadRequest, "your order_by_fields is wrong")
		return
	}
	if len(orderBy) > 0 {
		dbReq = dbReq.Order(orderBy...)
	}

	dbReq = dbReq.Limit(int(request.Limit.Value)).Offset(int(request.Offset.Value))

	err = dbReq.Select()
	if err != nil {
		logger.Log.Error(err)
		AnswerError(c, http.StatusInternalServerError, "some shit happened, call the admin")
		return
	}

	AnswerResponse(c, results)
}

func orderByFieldsToGoPg(value string, typeOfResult reflect.Type) ([]string, error) {
	results := []string{}

	fields := strings.Split(value, ",")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		match, err := regexp.MatchString("^-?[_a-zA-Z0-9]{1,128}$", field)
		if err != nil {
			return nil, err
		}
		if !match {
			return nil, errors.Errorf("\"%v\" doesn't match", field)
		}

		// TODO: restrict to tags

		if strings.HasPrefix(field, "-") {
			field = field[1:]
			results = append(results, field+" DESC")
		} else {
			results = append(results, field+" ASC")
		}
	}
	return results, nil
}
