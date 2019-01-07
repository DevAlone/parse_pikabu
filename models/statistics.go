package models

type NumberOfUsersToProcessEntry struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"ordering,filter"`
	Value     int64         `json:"value"`
}

func init() {
	for _, item := range []interface{}{
		&NumberOfUsersToProcessEntry{},
	} {
		Tables = append(Tables, item)
	}
}
