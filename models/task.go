package models

// Task is base task class
type Task struct {
	Id             uint64        `json:"id"`
	IsDone         bool          `sql:",notnull,default:false" json:"is_done"`
	IsTaken        bool          `sql:",notnull,default:false" json:"is_taken"`
	AddedTimestamp TimestampType `sql:",notnull" json:"added_timestamp"`
}

// SimpleTask is for tasks like "to parse communities"
type SimpleTask struct {
	Task
	Name string `sql:",notnull"`
}

/*
type ParseUserByUsernameTask struct {
	Task
	Username string `sql:",notnull" json:"username"`
}

type ParseUserByIdTask struct {
	Task
	PikabuId uint64 `sql:",notnull" json:"pikabu_id"`
}
*/

// ParseUserTask is a task for parsing users
type ParseUserTask struct {
	PikabuId       uint64        `sql:",pk" json:"pikabu_id"`
	IsDone         bool          `sql:",notnull,default:false" json:"is_done"`
	IsTaken        bool          `sql:",notnull,default:false" json:"is_taken"`
	AddedTimestamp TimestampType `sql:",notnull" json:"added_timestamp"`
	Username       string        `sql:",notnull" json:"username"`
}

// ParseStoryTask is a task for parsing stories by id
type ParseStoryTask struct {
	PikabuID       uint64        `sql:",pk" json:"pikabu_id"`
	AddedTimestamp TimestampType `sql:",notnull" json:"added_timestamp"`
}

func init() {
	Tables = append(Tables, []interface{}{
		&SimpleTask{},
		&ParseUserTask{},
	})

	/*
		CustomQueries = append(CustomQueries, `
			CREATE MATERIALIZED VIEW IF NOT EXISTS parse_user_by_username_tasks_is_not_done_and_is_not_taken
			AS SELECT * FROM parse_user_by_username_tasks WHERE is_done = false AND is_taken = false;
		`)
		CustomQueries = append(CustomQueries, `
			CREATE MATERIALIZED VIEW IF NOT EXISTS parse_user_by_id_tasks_is_not_done_and_is_not_taken
			AS SELECT * FROM parse_user_by_id_tasks WHERE is_done = false AND is_taken = false;
		`)
		CustomQueries = append(CustomQueries, `
			CREATE MATERIALIZED VIEW IF NOT EXISTS simple_tasks_is_not_done_and_is_not_taken
			AS SELECT * FROM simple_tasks WHERE is_done = false AND is_taken = false;
		`)
	*/

	addIndex("simple_tasks", "is_done", "")
	addIndex("simple_tasks", "is_taken", "")
	addIndex("simple_tasks", "added_timestamp", "")
	addIndex("simple_tasks", "name", "hash")

	addIndex("parse_user_tasks", "is_done", "")
	addIndex("parse_user_tasks", "is_taken", "")
	addIndex("parse_user_tasks", "added_timestamp", "")
	addIndex("parse_user_tasks", "username", "hash")
	addIndex("parse_user_tasks", "LOWER(username)", "hash")
}
