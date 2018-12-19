package main

import (
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/parser"
	"fmt"
	"github.com/go-pg/pg/orm"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		_, err := os.Stderr.WriteString(fmt.Sprintf(`Please, specify a command.
Available commands are: 
- core
- parser
`))
		if err != nil {
			panic(err)
		}
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
		if err != nil {
			panic(err)
		}

		// clear tables
		for _, table := range models.Tables {
			err := models.Db.DropTable(table, &orm.DropTableOptions{
				IfExists: true,
				Cascade:  true,
			})
			if err != nil {
				panic(err)
			}
		}
	default:
		_, err := os.Stderr.WriteString(fmt.Sprintf("Unknown command: %v", command))
		if err != nil {
			panic(err)
		}
	}
}
