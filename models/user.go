package models

type PikabuUser struct {
	PikabuId uint64 `sql:",pk" json:"pikabu_id" api:"ordering,filter"`

	Username            string        `sql:",notnull" gen_versions:"" json:"username" api:"ordering,filter"`
	Gender              string        `sql:",notnull" gen_versions:"" json:"gender" api:"ordering"`
	Rating              int32         `sql:",notnull" gen_versions:"" json:"rating" api:"ordering,filter"`
	NumberOfComments    int32         `sql:",notnull" gen_versions:"" json:"number_of_comments" api:"ordering,filter"`
	NumberOfSubscribers int32         `sql:",notnull" gen_versions:"" json:"number_of_subscribers" api:"ordering,filter"`
	NumberOfStories     int32         `sql:",notnull" gen_versions:"" json:"number_of_stories" api:"ordering"`
	NumberOfHotStories  int32         `sql:",notnull" gen_versions:"" json:"number_of_hot_stories" api:"ordering"`
	NumberOfPluses      int32         `sql:",notnull" gen_versions:"" json:"number_of_pluses" api:"ordering"`
	NumberOfMinuses     int32         `sql:",notnull" gen_versions:"" json:"number_of_minuses" api:"ordering"`
	SignupTimestamp     TimestampType `sql:",notnull" gen_versions:"" json:"signup_timestamp" api:"ordering"`
	AvatarURL           string        `sql:",notnull" gen_versions:"" json:"avatar_url" api:"ordering"`
	ApprovedText        string        `sql:",notnull" gen_versions:"" json:"approved_text" api:"ordering"`
	AwardIds            []uint64      `sql:",notnull,array" gen_versions:"" json:"award_ids" api:"ordering"`
	CommunityIds        []uint64      `sql:",notnull,array" gen_versions:"" json:"community_ids" api:"ordering"`
	BanHistoryItemIds   []uint64      `sql:",notnull,array" gen_versions:"" json:"ban_history_item_ids" api:"ordering"`
	BanEndTimestamp     TimestampType `sql:",notnull" gen_versions:"" json:"ban_end_timestamp" api:"ordering"`
	IsRatingHidden      bool          `sql:",notnull" gen_versions:"" json:"is_rating_hidden" api:"ordering"`
	IsBanned            bool          `sql:",notnull" gen_versions:"" json:"is_banned" api:"ordering"`
	IsPermanentlyBanned bool          `sql:",notnull" gen_versions:"" json:"is_permanently_banned" api:"ordering"`

	// ?
	// IsDeleted bool `sql:",notnull,default:false"`

	AddedTimestamp      TimestampType `sql:",notnull" json:"added_timestamp" api:"ordering"`
	LastUpdateTimestamp TimestampType `sql:",notnull" json:"last_update_timestamp" api:"ordering"`
	NextUpdateTimestamp TimestampType `sql:",notnull" json:"next_update_timestamp" api:"ordering"`
}

type PikabuUserAward struct {
	PikabuId uint64 `sql:",pk"`

	UserId uint64 `sql:",notnull" gen_versions:""`
	// TODO: figure out what the heck it is,
	// l4rever has 0 in this field in one of his awards
	AwardId       uint64 `sql:",notnull" gen_versions:""`
	AwardTitle    string `sql:",notnull" gen_versions:""`
	AwardImageURL string `sql:",notnull" gen_versions:""`
	StoryId       uint64 `sql:",notnull" gen_versions:""`
	StoryTitle    string `sql:",notnull" gen_versions:""`
	IssuingDate   string `sql:",notnull" gen_versions:""`
	// TODO: replace to bool?
	IsHidden  bool   `sql:",notnull" gen_versions:""`
	CommentId uint64 `sql:",notnull" gen_versions:""`
	// link to reason of award whether it was comment, story or anything else
	Link string `sql:",notnull" gen_versions:""`

	AddedTimestamp      TimestampType `sql:",notnull"`
	LastUpdateTimestamp TimestampType `sql:",notnull"`
}

type PikabuUserCommunity struct {
	Id uint64

	Name      string `sql:",notnull"`
	Link      string `sql:",notnull,unique"`
	AvatarURL string `sql:",notnull"`

	AddedTimestamp      TimestampType `sql:",notnull"`
	LastUpdateTimestamp TimestampType `sql:",notnull"`
}

type PikabuUserBanHistoryItem struct {
	PikabuId uint64 `sql:",pk"`

	BanStartTimestamp TimestampType `sql:",notnull" gen_versions:""`
	// id of comment caused ban if there was such
	CommentId               uint64 `sql:",notnull" gen_versions:""`
	CommentHtmlDeleteReason string `sql:",notnull" gen_versions:""`
	StoryId                 uint64 `sql:",notnull" gen_versions:""`
	UserId                  uint64 `sql:",notnull" gen_versions:""`
	BanReason               string `sql:",notnull" gen_versions:""`
	BanReasonId             uint64 `sql:",notnull" gen_versions:""`
	StoryURL                string `sql:",notnull" gen_versions:""`
	ModeratorId             uint64 `sql:",notnull" gen_versions:""`
	ModeratorName           string `sql:",notnull" gen_versions:""`
	ModeratorAvatar         string `sql:",notnull" gen_versions:""`
	// TODO: figure out what it means
	ReasonsLimit uint64 `sql:",notnull" gen_versions:""`
	ReasonCount  uint64 `sql:",notnull" gen_versions:""`
	ReasonTitle  string `sql:",notnull" gen_versions:""`

	AddedTimestamp      TimestampType `sql:",notnull"`
	LastUpdateTimestamp TimestampType `sql:",notnull"`
}

func init() {
	for _, item := range []interface{}{
		&PikabuUser{},
		&PikabuUserAward{},
		&PikabuUserCommunity{},
		&PikabuUserBanHistoryItem{},
	} {
		Tables = append(Tables, item)
	}

	/* // make it working
	CustomQueries = append(CustomQueries, `
		CREATE EXTENSION pg_trgm;
	`)
	*/

	// TODO: consider two indices for username(btree and hash)
	addIndex("pikabu_users", "username", "")
	addIndex("pikabu_users", "LOWER(username)", "hash")
	addIndex("pikabu_users", "username gin_trgm_ops", "gin")
	addIndex("pikabu_users", "added_timestamp", "")
	addIndex("pikabu_users", "last_update_timestamp", "")
	addIndex("pikabu_users", "next_update_timestamp", "")

	/*addUniqueIndex("core_user", "username", "")
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
	addIndex("core_userminusescountentry", "timestamp", "")*/
}
