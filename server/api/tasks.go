package api

import (
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"math/rand"
	"net/http"
	"strings"
)

func GetAnyTask(c *gin.Context) {
	taskData, err := TryToGetTaskFromDb()
	if err != nil {
		panic(err)
	}

	if taskData != nil {
		c.JSON(
			http.StatusOK,
			taskData,
		)
		return
	}

	c.JSON(http.StatusNotFound, map[string]string{
		"status":        "error",
		"error_message": "task not found",
	})
}

func TaskToResult(taskName string, taskData interface{}) interface{} {
	var result struct {
		Name string      `json:"name"`
		Data interface{} `json:"data"`
	}

	if strings.HasSuffix(taskName, "_tasks") {
		taskName = taskName[0 : len(taskName)-6]
	}

	result.Name = taskName
	result.Data = taskData

	return result
}

func TryToGetTaskFromDb() (interface{}, error) {
	tables := []helpers.Tuple{
		{"parse_user_by_username_tasks", &models.ParseUserByUsernameTask{}},
		{"parse_user_by_id_tasks", &models.ParseUserByIdTask{}},
	}
	rand.Shuffle(len(tables), func(i, j int) {
		tables[i], tables[j] = tables[j], tables[i]
	})

	for _, table := range tables {
		var result = table.Right

		err := models.Db.Model(result).
			Where("is_done = false AND is_taken = false	").
			OrderExpr("random()").
			Limit(1).
			Select(result)

		if err == pg.ErrNoRows {
			continue
		} else if err != nil {
			return nil, err
		}

		return TaskToResult(table.Left.(string), result), nil
	}

	return nil, nil
}
