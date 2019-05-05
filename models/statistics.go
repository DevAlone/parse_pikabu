package models

// NumberOfUsersToProcessEntry - number of users in queue
type NumberOfUsersToProcessEntry struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfStoriesToProcessEntry - number of stories in queue
type NumberOfStoriesToProcessEntry struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfCommentsToProcessEntry - number of comments in queue
type NumberOfCommentsToProcessEntry struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfTasksInChannelUpdateUserTask - number of task in channel
type NumberOfTasksInChannelUpdateUserTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfTasksInChannelParseNewUserTask - number of task in channel
type NumberOfTasksInChannelParseNewUserTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfTasksInChannelParseDeletedOrNeverExistedUserTask  -
type NumberOfTasksInChannelParseDeletedOrNeverExistedUserTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfTasksInChannelUpdateStoryTask -
type NumberOfTasksInChannelUpdateStoryTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfTasksInChannelParseNewStoryTask -
type NumberOfTasksInChannelParseNewStoryTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfTasksInChannelParseDeletedOrNeverExistedStoryTask -
type NumberOfTasksInChannelParseDeletedOrNeverExistedStoryTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfTasksInChannelParseAllCommunitiesTask -
type NumberOfTasksInChannelParseAllCommunitiesTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

func init() {
	Tables = append(Tables, []interface{}{
		&NumberOfUsersToProcessEntry{},
		&NumberOfStoriesToProcessEntry{},
		&NumberOfCommentsToProcessEntry{},
		&NumberOfTasksInChannelUpdateUserTask{},
		&NumberOfTasksInChannelParseNewUserTask{},
		&NumberOfTasksInChannelParseDeletedOrNeverExistedUserTask{},
		&NumberOfTasksInChannelUpdateStoryTask{},
		&NumberOfTasksInChannelParseNewStoryTask{},
		&NumberOfTasksInChannelParseDeletedOrNeverExistedStoryTask{},
		&NumberOfTasksInChannelParseAllCommunitiesTask{},
	}...)
}
