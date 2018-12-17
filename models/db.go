package models

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/logging"
	"errors"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"strings"
)

type TimestampType int32

var Db *pg.DB
var Tables []interface{}
var createIndexQueries []string

type QueryHook struct{}

func (this QueryHook) BeforeQuery(event *pg.QueryEvent) {}
func (this QueryHook) AfterQuery(event *pg.QueryEvent) {
	// TODO: make sure that it works
	query, err := event.FormattedQuery()
	if err != nil {
		logging.Log.Critical(err)
		return
	}
	if config.Settings.Debug {
		// logging.Log.Debug(fmt.Sprintf("SQL: %s %s\n", time.Since(event.Time), query))
		logging.Log.Debug(fmt.Sprintf("SQL: %s\n", query))
	}
}

func InitDb() error {
	logging.Log.Info("start initializing database")

	// TODO:
	/*&Community{},
	&CommunityCountersEntry{},

	&Image{},

	&Comment{},
	&CommentImagesVersion{},
	&CommentParentIdVersion{},
	&CommentCreatingTimestampVersion{},
	&CommentRatingVersion{},
	&CommentStoryIdVersion{},
	&CommentUserIdVersion{},
	&CommentAuthorUsernameVersion{},
	&CommentIsHiddenVersion{},
	&CommentIsDeletedVersion{},
	&CommentIsAuthorCommunityModeratorVersion{},
	&CommentIsAuthorPikabuTeamVersion{},
	&CommentTextVersion{},

	&StatisticsUsersInQueueCount{},*/

	dbConfig := config.Settings.Database
	Db = pg.Connect(&pg.Options{
		Database: dbConfig["Name"],
		User:     dbConfig["Username"],
		Password: dbConfig["Password"],
	})

	Db.AddQueryHook(QueryHook{})

	logging.Log.Info("creating schema")
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
		logging.Log.Info("creating index:", query)
		_, err = Db.Exec(query)
		if err != nil {
			return err
		}
	}

	return err
}

func createSchema() error {
	logging.Log.Debugf("number of Tables is %d", len(Tables))
	for i, model := range Tables {
		err := Db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			print(i, model, "\n")
			panic(err)
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

	index := "CREATE "
	if unique {
		index += "UNIQUE "
	}
	index += "INDEX IF NOT EXISTS " + tableName + "__"
	index += strings.Join(columns, "__")
	if unique {
		index += "__unique"
	}
	index += "__idx ON " + tableName

	if len(method) > 0 {
		index += " USING " + method + " "
	}

	index += "("
	index += strings.Join(columns, ",")
	index += ");\n"
	createIndexQueries = append(createIndexQueries, index)
}

type FieldVersionBase struct {
	Timestamp TimestampType `sql:",pk,notnull"`
	ItemId    uint64        `sql:",pk,notnull"`
}

type Int32FieldVersion struct {
	FieldVersionBase
	Value int32 `sql:",notnull"`
}
type Int64FieldVersion struct {
	FieldVersionBase
	Value int64 `sql:",notnull"`
}
type UInt32FieldVersion struct {
	FieldVersionBase
	Value uint32 `sql:",notnull"`
}
type UInt64FieldVersion struct {
	FieldVersionBase
	Value uint64 `sql:",notnull"`
}
type StringFieldVersion struct {
	FieldVersionBase
	Value string `sql:",notnull"`
}
type BoolFieldVersion struct {
	FieldVersionBase
	Value bool `sql:",notnull"`
}
type TimestampTypeFieldVersion struct {
	FieldVersionBase
	Value TimestampType `sql:",notnull"`
}
