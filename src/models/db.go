package models

import (
	. "config"
	"errors"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"logging"
	"strings"
	"time"
)

var Db *pg.DB
var tables []interface{}
var createIndicesQuery string = ""

func InitDb() error {
	logging.Log.Info("start initializing database")

	tables = []interface{}{
		&User{},
		&UserRatingEntry{},
		&UserSubscribersCountEntry{},
		&UserCommentsCountEntry{},
		&UserPostsCountEntry{},
		&UserHotPostsCountEntry{},
		&UserPlusesCountEntry{},
		&UserMinusesCountEntry{},
		&UserAvatarURLVersion{},
		&PikabuUser{},

		&Community{},
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

		&StatisticsUsersInQueueCount{},
	}

	dbConfig := Settings.Database
	Db = pg.Connect(&pg.Options{
		Database: dbConfig["Name"],
		User:     dbConfig["Username"],
		Password: dbConfig["Password"],
	})

	Db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			logging.Log.Critical(err)
			return
		}
		if Settings.Debug {
			logging.Log.Debug(fmt.Sprintf("SQL: %s %s\n", time.Since(event.StartTime), query))
		}
	})

	logging.Log.Info("creating schema")
	err := createSchema()
	if err != nil {
		return err
	}

	for _, query := range strings.Split(createIndicesQuery, ";") {
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
	for _, model := range tables {
		err := Db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
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
	createIndicesQuery += index
}

type FieldVersionBase struct {
	Timestamp int32  `sql:",pk,notnull"`
	ItemId    uint64 `sql:",pk,notnull"`
}

type Int32FieldVersion struct {
	FieldVersionBase
	Value int32 `sql:",notnull"`
}
type Int64FieldVersion struct {
	FieldVersionBase
	Value int64 `sql:",notnull"`
}
type Uint64FieldVersion struct {
	FieldVersionBase
	Value uint64 `sql:",notnull"`
}
type TextFieldVersion struct {
	FieldVersionBase
	Value string `sql:",notnull"`
}
type BoolFieldVersion struct {
	FieldVersionBase
	Value bool `sql:",notnull"`
}
