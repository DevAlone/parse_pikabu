package models

type Task struct {
	Id             uint64        `json:"id"`
	IsDone         bool          `sql:",notnull,default:false" json:"is_done"`
	IsTaken        bool          `sql:",notnull,default:false" json:"is_taken"`
	AddedTimestamp TimestampType `sql:",notnull" json:"added_timestamp"`
}

type SimpleTask struct {
	Task
	Name string `sql:",notnull"`
}

type ParseUserByUsernameTask struct {
	Task
	Username string `sql:",notnull" json:"username"`
}

type ParseUserByIdTask struct {
	Task
	PikabuId uint64 `sql:",notnull" json:"pikabu_id"`
}

func init() {
	for _, table := range []interface{}{
		&SimpleTask{},
		&ParseUserByUsernameTask{},
		&ParseUserByIdTask{},
	} {
		Tables = append(Tables, table)
	}

	addIndex("simple_tasks", "is_done", "")
	addIndex("simple_tasks", "is_taken", "")
	addIndex("simple_tasks", "added_timestamp", "")
	addIndex("simple_tasks", "name", "hash")

	addIndex("parse_user_by_username_tasks", "is_done", "")
	addIndex("parse_user_by_username_tasks", "is_taken", "")
	addIndex("parse_user_by_username_tasks", "added_timestamp", "")
	addIndex("parse_user_by_username_tasks", "username", "hash")

	addIndex("parse_user_by_id_tasks", "is_done", "")
	addIndex("parse_user_by_id_tasks", "is_taken", "")
	addIndex("parse_user_by_id_tasks", "added_timestamp", "")
	addIndex("parse_user_by_id_tasks", "pikabu_id", "")
}
