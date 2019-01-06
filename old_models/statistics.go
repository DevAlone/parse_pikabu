package old_models

type StatisticsUsersInQueueCount struct {
	Timestamp int64 `sql:",pk"`
	Value     int32
}
