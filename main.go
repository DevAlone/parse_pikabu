package main

import (
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/parser"
	"bitbucket.org/d3dev/parse_pikabu/server/middlewares"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg/orm"
	"os"
	"strings"
)

func handleError(err error) {
	if err == nil {
		return
	}
	/*if e, ok := err.(*errors.Error); ok {
		_, err := os.Stderr.WriteString(e.ErrorStack())
		if err != nil {
			panic(err)
		}
		_, err = os.Stderr.WriteString(e.Error())
		if err != nil {
			panic(err)
		}
	}*/

	if e, ok := err.(*errors.Error); ok {
		_, er := os.Stderr.WriteString(e.ErrorStack())
		if er != nil {
			panic(er)
		}
	}
	panic(err)
}

func main() {
	if len(os.Args) < 2 {
		_, err := os.Stderr.WriteString(fmt.Sprintf(`Please, specify a command.
Available commands are: 
- core
- parser
- clean_db
- add_parser
`))
		handleError(err)
		return
	}

	command := os.Args[1]
	os.Args = os.Args[1:]

	switch command {
	case "core":
		Main()
	case "parser":
		parser.Main()
	case "clean_db":
		err := models.InitDb()
		handleError(err)

		// clear tables
		for _, table := range models.Tables {
			err := models.Db.DropTable(table, &orm.DropTableOptions{
				IfExists: true,
				Cascade:  true,
			})
			handleError(err)
		}
	case "add_parser":
		redisClient := helpers.GetRedisClient()
		// parse_pikabu_server_authentication_middleware_session_group_
		//
		if len(os.Args) < 2 {
			handleError(errors.New("too few arguments"))
		}
		key := "parse_pikabu_server_authentication_middleware_session_group_" + strings.TrimSpace(os.Args[1])
		err := redisClient.Set(key, fmt.Sprint(middlewares.GROUP_PARSER), 0).Err()
		handleError(err)
	case "load_from_old_db":
		loadFromOldDb()
	default:
		_, err := os.Stderr.WriteString(fmt.Sprintf("Unknown command: %v", command))
		if err != nil {
			handleError(err)
		}
	}
}
