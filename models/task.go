package models

type Task struct {
	Id             uint64        `json:"-"`
	IsDone         bool          `sql:",notnull,default:false" json:"is_done"`
	IsTaken        bool          `sql:",notnull,default:false" json:"is_taken"`
	AddedTimestamp TimestampType `sql:",notnull" json:"added_timestamp"`
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
		&ParseUserByUsernameTask{},
		&ParseUserByIdTask{},
	} {
		tables = append(tables, table)
	}
}
