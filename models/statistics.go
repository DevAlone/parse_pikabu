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

// NumberOfItemsInChannelUpdateUserTask - number of task in channel
type NumberOfItemsInChannelUpdateUserTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfItemsInChannelParseNewUserTask - number of task in channel
type NumberOfItemsInChannelParseNewUserTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfItemsInChannelParseDeletedOrNeverExistedUserTask  -
type NumberOfItemsInChannelParseDeletedOrNeverExistedUserTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfItemsInChannelUpdateStoryTask -
type NumberOfItemsInChannelUpdateStoryTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfItemsInChannelParseNewStoryTask -
type NumberOfItemsInChannelParseNewStoryTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfItemsInChannelParseDeletedOrNeverExistedStoryTask -
type NumberOfItemsInChannelParseDeletedOrNeverExistedStoryTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfItemsInChannelParseAllCommunitiesTask -
type NumberOfItemsInChannelParseAllCommunitiesTask struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

// NumberOfItemsInChannelParserResults -
type NumberOfItemsInChannelParserResults struct {
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"order,filter"`
	Value     int64         `json:"value"`
}

func init() {
	Tables = append(Tables, []interface{}{
		&NumberOfUsersToProcessEntry{},
		&NumberOfStoriesToProcessEntry{},
		&NumberOfCommentsToProcessEntry{},
		&NumberOfItemsInChannelUpdateUserTask{},
		&NumberOfItemsInChannelParseNewUserTask{},
		&NumberOfItemsInChannelParseDeletedOrNeverExistedUserTask{},
		&NumberOfItemsInChannelUpdateStoryTask{},
		&NumberOfItemsInChannelParseNewStoryTask{},
		&NumberOfItemsInChannelParseDeletedOrNeverExistedStoryTask{},
		&NumberOfItemsInChannelParseAllCommunitiesTask{},

		&NumberOfItemsInChannelParserResults{},
	}...)
}
