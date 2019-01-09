package main

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/parser"
	"bitbucket.org/d3dev/parse_pikabu/server/middlewares"
	"github.com/go-errors/errors"
)

func addParsers() error {
	redisClient := helpers.GetRedisClient()
	if len(os.Args) < 2 {
		handleError(errors.New("too few arguments"))
	}
	parsersConfig, err := parser.NewParsersConfigFromFile(os.Args[1])
	if err != nil {
		return err
	}

	for _, parserConfig := range parsersConfig.Configs {
		sessionId := strings.TrimSpace(parserConfig.ApiSessionId)
		key := "parse_pikabu_server_authentication_middleware_session_group_" + sessionId
		err := redisClient.Set(key, fmt.Sprint(middlewares.GROUP_PARSER), 0).Err()

		if err != nil {
			return err
		}
	}

	return nil
}
