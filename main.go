package main

import (
	"bitbucket.org/d3dev/parse_pikabu/amqp_helper"
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/d3dev/parse_pikabu/core"
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/server/middlewares"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/parser"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg/orm"
)

var commands = map[string]func(){
	"core": func() {
		core.Main()
	},
	"parser": func() {
		parser.Main()
	},
	"clean_db": func() {
		err := models.InitDb()
		if err != nil {
			helpers.PanicOnError(err)
		}

		// clear tables
		for _, table := range models.Tables {
			err := models.Db.DropTable(table, &orm.DropTableOptions{
				IfExists: true,
				Cascade:  true,
			})
			if err != nil {
				helpers.PanicOnError(err)
			}
		}
	},
	"add_parser": func() {
		redisClient := helpers.GetRedisClient()
		if len(os.Args) < 2 {
			helpers.PanicOnError(errors.New("too few arguments"))
		}
		key := "parse_pikabu_server_authentication_middleware_session_group_" + strings.TrimSpace(os.Args[1])
		err := redisClient.Set(key, fmt.Sprint(middlewares.GROUP_PARSER), 0).Err()
		helpers.PanicOnError(err)
	},
	"add_parsers_from_config": func() {
		err := addParsersFromConfig()
		helpers.PanicOnError(err)
	},
	"load_from_old_db": func() {
		loadFromOldDb()
	},
}

func main() {
	if len(os.Args) < 2 {
		var commandsList string
		for command := range commands {
			commandsList += "\t-" + command + "\n"
		}
		_, err := os.Stderr.WriteString(fmt.Sprintf(`Please, specify a command.
Available commands are: 
%s
`, commandsList))
		helpers.PanicOnError(err)
		return
	}

	command := strings.TrimSpace(os.Args[1])
	os.Args = os.Args[1:]

	configFilePath := flag.String("config", "core.config.json", "config file")

	if configFilePath == nil || len(*configFilePath) == 0 {
		panic(errors.New("configFilePath is nil"))
	}

	err := config.UpdateSettingsFromFile(*configFilePath)
	if err != nil {
		helpers.PanicOnError(err)
	}

	if handler, found := commands[command]; found {
		handler()
	} else {
		helpers.PanicOnError(errors.Errorf("wrong command"))
	}

	helpers.PanicOnError(amqp_helper.Cleanup())
}
