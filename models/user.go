package models

// TODO: update structure

type PikabuUser struct {
	PikabuId uint64 `sql:",pk"`

	Username            string        `sql:",notnull" gen_versions:""`
	Gender              string        `sql:",notnull" gen_versions:""`
	Rating              float32       `sql:",notnull" gen_versions:""`
	NumberOfComments    int32         `sql:",notnull" gen_versions:""`
	NumberOfSubscribers int32         `sql:",notnull" gen_versions:""`
	NumberOfStories     int32         `sql:",notnull" gen_versions:""`
	NumberOfHotStories  int32         `sql:",notnull" gen_versions:""`
	NumberOfPluses      int32         `sql:",notnull" gen_versions:""`
	NumberOfMinuses     int32         `sql:",notnull" gen_versions:""`
	SignupTimestamp     TimestampType `sql:",notnull" gen_versions:""`
	AvatarURL           string        `sql:",notnull" gen_versions:""`
	ApprovedText        string        `sql:",notnull" gen_versions:""`
	AwardIds            []uint64      `sql:",notnull,array" gen_versions:""`
	CommunityIds        []uint64      `sql:",notnull,array" gen_versions:""`
	BanHistoryItemIds   []uint64      `sql:",notnull,array" gen_versions:""`
	BanEndTimestamp     TimestampType `sql:",notnull" gen_versions:""`
	IsRatingHidden      bool          `sql:",notnull" gen_versions:""`
	IsBanned            bool          `sql:",notnull" gen_versions:""`
	IsPermanentlyBanned bool          `sql:",notnull" gen_versions:""`

	// ?
	// IsDeleted bool `sql:",notnull,default:false"`

	AddedTimestamp TimestampType `sql:",notnull"`
	// PreviousUpdateTimestamp TimestampType `sql:",notnull, default:0"`
	LastUpdateTimestamp TimestampType `sql:",notnull,default:0"`
	NextUpdateTimestamp TimestampType `sql:",notnull,default:0"`
}

// TODO: condider generating versions for this struct
type PikabuUserAward struct {
	PikabuId uint64 `sql:",pk"`

	Timestamp TimestampType `sql:",notnull"`
	UserId    uint64        `sql:",notnull"`
	// TODO: figure out what the heck it is,
	// l4rever has 0 in this field in one of his awards
	AwardId       uint64 `sql:",notnull"`
	AwardTitle    string `sql:",notnull"`
	AwardImageURL string `sql:",notnull"`
	StoryId       uint64 `sql:",notnull"`
	StoryTitle    string `sql:",notnull"`
	IssuingDate   string `sql:",notnull"`
	// TODO: replace to bool?
	IsHidden  bool   `sql:",notnull"`
	CommentId uint64 `sql:",notnull"`
	// link to reason of award whether it was comment, story or anything else
	Link string `sql:",notnull"`
}

type PikabuUserCommunity struct {
	Id uint64

	Timestamp TimestampType `sql:",notnull"`
	Name      string        `sql:",notnull"`
	Link      string        `sql:",notnull,unique"`
	AvatarURL string        `sql:",notnull"`
}

type PikabuUserBanHistoryItem struct {
	PikabuId uint64 `sql:",pk"`

	Timestamp         TimestampType `sql:",notnull"`
	BanStartTimestamp TimestampType `sql:",notnull"`
	// id of comment caused ban if there was such
	CommentId               uint64 `sql:",notnull"`
	CommentHtmlDeleteReason string `sql:",notnull"`
	StoryId                 uint64 `sql:",notnull"`
	UserId                  uint64 `sql:",notnull"`
	BanReason               string `sql:",notnull"`
	BanReasonId             uint64 `sql:",notnull"`
	StoryURL                string `sql:",notnull"`
	ModeratorId             uint64 `sql:",notnull"`
	ModeratorName           string `sql:",notnull"`
	ModeratorAvatar         string `sql:",notnull"`
	// TODO: figure out what it means
	ReasonsLimit uint64 `sql:",notnull"`
	ReasonCount  uint64 `sql:",notnull"`
	ReasonTitle  string `sql:",notnull"`
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

	// TODO: add proper indices
	// username is not hash for fast sorting
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
	addIndex("core_userminusescountentry", "timestamp", "")*/
}
