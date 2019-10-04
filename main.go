package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"sync"

	"github.com/DevAlone/parse_pikabu/amqphelper"
	"github.com/DevAlone/parse_pikabu/core"
	"github.com/DevAlone/parse_pikabu/core/config"
	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/DevAlone/parse_pikabu/core/server/middlewares"
	"github.com/DevAlone/parse_pikabu/globals"
	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/DevAlone/parse_pikabu/parser"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/profile"
)

var commands = map[string]func(){
	"single_process_mode": func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			core.Main()
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			parser.Main()
			wg.Done()
		}()
		wg.Wait()
	},
	"core": func() {
		core.Main()
	},
	"parser": func() {
		// TODO: modify to fetch tasks from amqp or something
		parser.Main()
	},
	"clean_db": func() {
		logger.Init()

		err := models.InitDb()
		helpers.PanicOnError(err)

		// clear tables
		for _, table := range models.Tables {
			err := models.Db.DropTable(table, &orm.DropTableOptions{
				IfExists: true,
				Cascade:  true,
			})
			helpers.PanicOnError(err)
		}
	},
	"print_index_queries": func() {
		logger.Init()

		queries := models.GetIndexQueries()
		for _, query := range queries {
			fmt.Println(query)
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
	"fix_usernames": func() {
		fixUsernames()
	},
}

func main() {
	var err error

	if len(os.Args) < 2 {
		var commandsList string
		for command := range commands {
			commandsList += "\t-" + command + "\n"
		}
		_, err = os.Stderr.WriteString(fmt.Sprintf(`Please, specify a command.
Available commands are: 
%s
`, commandsList))
		helpers.PanicOnError(err)
		return
	}

	command := strings.TrimSpace(os.Args[1])
	os.Args = os.Args[1:]

	configFilePath := flag.String("config", "core.config.json", "config file")
	cpuProfile := flag.String("cpuprofile", "", "set to true to profile cpu")
	memProfile := flag.String("memprofile", "", "set to true to profile memory")
	doNotParseUsersFlag := flag.String("do-not-parse-users", "false", "do not parse users")
	doNotParseStoriesFlag := flag.String("do-not-parse-stories", "false", "do not parse stories")

	flag.Parse()

	globals.DoNotParseUsers = strings.HasPrefix(strings.ToLower(*doNotParseUsersFlag), "t")
	globals.DoNotParseStories = strings.HasPrefix(strings.ToLower(*doNotParseStoriesFlag), "t")

	if configFilePath == nil || len(*configFilePath) == 0 {
		panic(errors.New("configFilePath is nil"))
	}

	if *cpuProfile != "" {
		f, err := os.Create("cpu.pprof")
		helpers.PanicOnError(err)
		err = pprof.StartCPUProfile(f)
		helpers.PanicOnError(err)
		defer pprof.StopCPUProfile()
	}

	if *memProfile != "" {
		p := profile.Start(profile.MemProfile)
		defer p.Stop()
	}

	err = config.UpdateSettingsFromFile(*configFilePath)
	helpers.PanicOnError(err)

	err = globals.Init()
	helpers.PanicOnError(err)

	if handler, found := commands[command]; found {
		handler()
	} else {
		helpers.PanicOnError(errors.Errorf("wrong command"))
	}

	helpers.PanicOnError(amqphelper.Cleanup())
}
