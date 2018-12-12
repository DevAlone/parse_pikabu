package models

type StatisticsUsersInQueueCount struct {
	Timestamp int64 `sql:",pk"`
	Value     int32
}
