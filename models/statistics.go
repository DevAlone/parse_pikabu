package models

// TODO: update structure

type StatisticsUsersInQueueCount struct {
	Timestamp int64 `sql:",pk"`
	Value     int32
}
