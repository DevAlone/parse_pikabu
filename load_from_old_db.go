package main

import (
	"fmt"

	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func loadFromOldDb() {
	fmt.Println("loading from db pikabot_graphs...")

	err := models.InitDb()
	if err != nil {
		handleError(err)
	}

	// clear tables
	for _, table := range models.Tables {
		err := models.Db.DropTable(table, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			handleError(err)
		}
	}

	// create again
	err = models.InitDb()
	if err != nil {
		handleError(err)
	}

	oldDb := pg.Connect(&pg.Options{
		Database: "pikabot_graphs",
		User:     "pikabot_graphs",
		Password: "pikabot_graphs",
	})
}
