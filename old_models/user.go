package old_models

type User struct {
	TableName struct{} `sql:"core_user"`

	Id                  uint64
	Username            string `sql:",type:varchar(64),notnull"`
	Rating              int32  `sql:",notnull"`
	CommentsCount       int32  `sql:",notnull"`
	PostsCount          int32  `sql:",notnull"`
	HotPostsCount       int32  `sql:",notnull"`
	PlusesCount         int32  `sql:",notnull"`
	MinusesCount        int32  `sql:",notnull"`
	SubscribersCount    int32  `sql:",notnull"`
	IsRatingBan         bool   `sql:",notnull"`
	UpdatingPeriod      int32  `sql:",notnull"`
	AvatarURL           string
	Info                string
	IsUpdated           bool  `sql:",notnull"`
	LastUpdateTimestamp int64 `sql:",notnull"`
	Approved            string
	Awards              string
	Gender              string `sql:",notnull"`
	PikabuId            int64
	SignupTimestamp     int64 `sql:",notnull"`
	Deleted             bool  `sql:",notnull,default:false"`
}

type PikabuUser struct {
	TableName struct{} `sql:"core_pikabuuser"`

	PikabuId    int64  `sql:",pk"`
	Username    string `sql:",notnull"`
	IsProcessed bool   `sql:",notnull,default:false"`
}

type CountersEntryBase struct {
	Id        uint64
	UserId    int64
	Timestamp int64
	Value     int32
}

type UserRatingEntry struct {
	TableName struct{} `sql:"core_userratingentry"`
	CountersEntryBase
}
type UserSubscribersCountEntry struct {
	TableName struct{} `sql:"core_usersubscriberscountentry"`
	CountersEntryBase
}
type UserCommentsCountEntry struct {
	TableName struct{} `sql:"core_usercommentscountentry"`
	CountersEntryBase
}
type UserPostsCountEntry struct {
	TableName struct{} `sql:"core_userpostscountentry"`
	CountersEntryBase
}
type UserHotPostsCountEntry struct {
	TableName struct{} `sql:"core_userhotpostscountentry"`
	CountersEntryBase
}
type UserPlusesCountEntry struct {
	TableName struct{} `sql:"core_userplusescountentry"`
	CountersEntryBase
}
type UserMinusesCountEntry struct {
	TableName struct{} `sql:"core_userminusescountentry"`
	CountersEntryBase
}

// TODO: add comments count and other entries tables

type UserAvatarURLVersion struct {
	Timestamp int64  `sql:",pk,notnull"`
	ItemId    int64  `sql:",pk,notnull"`
	Value     string `sql:",notnull"`
}

func init() {
	/*

		// username is not hash for fast sorting
		addUniqueIndex("core_user", "username", "")
		addIndex("core_user", "rating", "")
		addIndex("core_user", "comments_count", "")
		addIndex("core_user", "posts_count", "")
		addIndex("core_user", "hot_posts_count", "")
		addIndex("core_user", "pluses_count", "")
		addIndex("core_user", "minuses_count", "")
		addIndex("core_user", "subscribers_count", "")
		addIndex("core_user", "is_rating_ban", "")
		addIndex("core_user", "updating_period", "")
		addIndex("core_user", "avatar_url", "hash")
		// addIndex("core_user", "info", "") // maybe some kind of substrings index
		addIndex("core_user", "is_updated", "")
		addIndex("core_user", "last_update_timestamp", "")
		// addIndex("core_user", "approved", "")
		// addIndex("core_user", "awards", "")
		addIndex("core_user", "gender", "hash")
		addIndex("core_user", "pikabu_id", "")
		addIndex("core_user", "signup_timestamp", "")
		addIndex("core_user", "deleted", "")

		addIndex("core_pikabuuser", "username", "hash")
		addIndex("core_pikabuuser", "is_processed", "")

		addUniqueIndex("core_userratingentry", []string{"user_id", "timestamp"}, "")
		addUniqueIndex("core_usersubscriberscountentry", []string{"user_id", "timestamp"}, "")
		addUniqueIndex("core_usercommentscountentry", []string{"user_id", "timestamp"}, "")
		addUniqueIndex("core_userpostscountentry", []string{"user_id", "timestamp"}, "")
		addUniqueIndex("core_userhotpostscountentry", []string{"user_id", "timestamp"}, "")
		addUniqueIndex("core_userplusescountentry", []string{"user_id", "timestamp"}, "")
		addUniqueIndex("core_userminusescountentry", []string{"user_id", "timestamp"}, "")

		addIndex("core_userratingentry", "user_id", "")
		addIndex("core_usersubscriberscountentry", "user_id", "")
		addIndex("core_usercommentscountentry", "user_id", "")
		addIndex("core_userpostscountentry", "user_id", "")
		addIndex("core_userhotpostscountentry", "user_id", "")
		addIndex("core_userplusescountentry", "user_id", "")
		addIndex("core_userminusescountentry", "user_id", "")

		addIndex("core_userratingentry", "timestamp", "")
		addIndex("core_usersubscriberscountentry", "timestamp", "")
		addIndex("core_usercommentscountentry", "timestamp", "")
		addIndex("core_userpostscountentry", "timestamp", "")
		addIndex("core_userhotpostscountentry", "timestamp", "")
		addIndex("core_userplusescountentry", "timestamp", "")
		addIndex("core_userminusescountentry", "timestamp", "")
	*/
}
