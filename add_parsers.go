package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/DevAlone/parse_pikabu/core/server/middlewares"
	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/parser"
	"github.com/go-errors/errors"
)

func addParsersFromConfig() error {
	redisClient := helpers.GetRedisClient()
	if len(os.Args) < 2 {
		return errors.New("too few arguments")
	}
	parsersConfig, err := parser.NewParsersConfigFromFile(os.Args[1])
	if err != nil {
		return err
	}

	for _, parserConfig := range parsersConfig.Configs {
		sessionId := strings.TrimSpace(parserConfig.APISessionID)
		key := "parse_pikabu_server_authentication_middleware_session_group_" + sessionId
		err := redisClient.Set(key, fmt.Sprint(middlewares.GROUP_PARSER), 0).Err()

		if err != nil {
			return err
		}
	}

	return nil
}
