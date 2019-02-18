package models

import (
	"errors"
	"strings"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

//go:generate bash -c "cd ..; go run models_versions_fields.generator.go; go fmt models/models.generated.go"

type TimestampType int64

var Db *pg.DB
var Tables []interface{}
var CustomQueries []string
var createIndexQueries []string

type QueryHook struct{}

func (this QueryHook) BeforeQuery(event *pg.QueryEvent) {}
func (this QueryHook) AfterQuery(event *pg.QueryEvent) {
	// TODO: make option in settings
	/*
		if config.Settings.Debug {
			query, err := event.FormattedQuery()
			if err != nil {
				logger.Log.Critical(err)
				return
			}
			// logger.Log.Debug(fmt.Sprintf("SQL: %s %s\n", time.Since(event.Time), query))
			logger.Log.Debug(fmt.Sprintf("SQL: %s\n", query))
		}
	*/
}

func InitDb() error {
	logger.Log.Info("start initializing database")

	dbConfig := config.Settings.Database
	Db = pg.Connect(&pg.Options{
		Database: dbConfig["Name"],
		User:     dbConfig["Username"],
		Password: dbConfig["Password"],
	})

	Db.AddQueryHook(QueryHook{})

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
			return err
		}
	}

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
			return err
		}
	}

	for i, query := range CustomQueries {
		_, err := Db.Exec(query)
		if err != nil {
			print(i, query, "\n")
			return err
		}
	}

	return nil
}

func addIndex(tableName string, _columns interface{}, method string) {
	_addIndex(tableName, _columns, method, false)
}

func addUniqueIndex(tableName string, _columns interface{}, method string) {
	_addIndex(tableName, _columns, method, true)
}

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
