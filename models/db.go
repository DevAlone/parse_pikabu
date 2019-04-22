package models

import (
	"fmt"
	"strings"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

//go:generate bash -c "cd ..; go run models_versions_fields.generator.go; go fmt models/models_versions_fields.generated.go"
//go:generate bash -c "cd ..; go run models_distributions_fields.generator.go; go fmt models/models_distributions_fields.generated.go"

// TimestampType - type for timestamp fields
type TimestampType int64

// Db is a database of project
var Db *pg.DB

// Tables - tables in db
var Tables []interface{}

// CustomQueries - use for custom queries like creating materiazed view
var CustomQueries []string
var createIndexQueries []string

type queryHook struct{}

func (qh queryHook) BeforeQuery(event *pg.QueryEvent) {}
func (qh queryHook) AfterQuery(event *pg.QueryEvent) {
	if config.Settings.Debug && config.Settings.LogSQLQueries {
		query, err := event.FormattedQuery()
		if err != nil {
			logger.Log.Fatal(err)
			return
		}
		// logger.Log.Debug(fmt.Sprintf("SQL: %s %s\n", time.Since(event.Time), query))
		logger.Log.Debug(fmt.Sprintf("SQL: %s\n", query))
	}
}

// InitDb initializes the database
func InitDb() error {
	logger.Log.Info("start initializing database")

	dbConfig := config.Settings.Database
	Db = pg.Connect(&pg.Options{
		Database: dbConfig["Name"],
		User:     dbConfig["Username"],
		Password: dbConfig["Password"],
	})

	Db.AddQueryHook(queryHook{})

	logger.Log.Info("creating schema")
	err := createSchema()
	if err != nil {
		return err
	}

	for _, query := range createIndexQueries {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}
		query += ";"
		logger.Log.Info("creating index:", query)
		_, err = Db.Exec(query)
		if err != nil {
			print("shit during creating index ", query, "\n")
			return errors.New(err)
		}
	}

	logger.Log.Info("Database created successfully")

	return err
}

func createSchema() error {
	logger.Log.Debugf("number of Tables is %d", len(Tables))
	for i, model := range Tables {
		err := Db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			print(i, model, "\n")
			return errors.New(err)
		}
	}
	logger.Log.Debugf("number of CustomQueries is %d", len(CustomQueries))

	for i, query := range CustomQueries {
		_, err := Db.Exec(query)
		if err != nil {
			print(i, query, "\n")
			return errors.New(err)
		}
	}

	return nil
}

func addIndex(tableName string, _columns interface{}, method string) {
	_addIndex(tableName, _columns, method, false)
}

/*
func addUniqueIndex(tableName string, _columns interface{}, method string) {
	_addIndex(tableName, _columns, method, true)
}
*/

func _addIndex(tableName string, _columns interface{}, method string, unique bool) {
	columns := []string{}

	switch val := _columns.(type) {
	case string:
		columns = append(columns, val)
	case []string:
		columns = val
	default:
		panic(errors.New("addIndex() bad 2nd argument"))
	}

	method = strings.ToLower(method)

	indexName := ""
	if unique {
		indexName += "unique__"
	}
	indexName += tableName + "__"
	if len(method) > 0 {
		indexName += "using_method_" + method + "__"
	}
	indexName += strings.Join(columns, "__") + "__"
	indexName += "idx"

	indexName = strings.Replace(indexName, "(", "_oparenthesis_", -1)
	indexName = strings.Replace(indexName, ")", "_cparenthesis_", -1)
	indexName = strings.Replace(indexName, " ", "_space_", -1)
	indexName = strings.ToLower(indexName)

	indexQuery := "CREATE "
	if unique {
		indexQuery += "UNIQUE "
	}
	indexQuery += "INDEX IF NOT EXISTS " + indexName + " "
	indexQuery += "ON " + tableName + " "
	if len(method) > 0 {
		indexQuery += "USING " + method + " "
	}

	indexQuery += "("
	indexQuery += strings.Join(columns, ",")
	indexQuery += ");\n"
	createIndexQueries = append(createIndexQueries, indexQuery)
}
