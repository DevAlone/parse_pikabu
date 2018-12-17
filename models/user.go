package models

// TODO: update structure

type PikabuUser struct {
	PikabuId uint64 `sql:",pk"`

	Username            string        `sql:",notnull"`
	Gender              string        `sql:",notnull"`
	Rating              float32       `sql:",notnull"`
	NumberOfComments    int32         `sql:",notnull"`
	NumberOfSubscribers int32         `sql:",notnull"`
	NumberOfStories     int32         `sql:",notnull"`
	NumberOfHotStories  int32         `sql:",notnull"`
	NumberOfPluses      int32         `sql:",notnull"`
	NumberOfMinuses     int32         `sql:",notnull"`
	SignupTimestamp     TimestampType `sql:",notnull"`
	AvatarURL           string        `sql:",notnull"`
	ApprovedText        string        `sql:",notnull"`
	Awards              []uint64      `sql:",notnull,array"`
	Communities         []uint64      `sql:",notnull,array"`
	BanHistory          []uint64      `sql:",notnull,array"`
	BanEndTimestamp     TimestampType `sql:",notnull"`
	IsRatingHidden      bool          `sql:",notnull"`
	IsBanned            bool          `sql:",notnull"`
	IsPermanentlyBanned bool          `sql:",notnull"`

	// ?
	IsDeleted bool `sql:",notnull,default:false"`

	AddedTimestamp TimestampType `sql:",notnull"`
	// PreviousUpdateTimestamp TimestampType `sql:",notnull, default:0"`
	LastUpdateTimestamp TimestampType `sql:",notnull,default:0"`
	NextUpdateTimestamp TimestampType `sql:",notnull,default:0"`
}

type PikabuUserAward struct {
	PikabuId uint64 `sql:",pk"`
	UserId   uint64 `sql:",notnull"`
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

	Name      string `sql:",notnull"`
	Link      string `sql:",notnull"`
	AvatarURL string `sql:",notnull"`
}

type PikabuUserBanHistoryItem struct {
	PikabuId uint64 `sql:",pk"`

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

type PikabuUserUsernameVersion struct{ StringFieldVersion }
type PikabuUserRatingVersion struct{ Float32FieldVersion }
type PikabuUserGenderVersion struct{ StringFieldVersion }
type PikabuUserNumberOfCommentsVersion struct{ UInt32FieldVersion }
type PikabuUserNumberOfStoriesVersion struct{ UInt32FieldVersion }
type PikabuUserNumberOfHotStoriesVersion struct{ UInt32FieldVersion }
type PikabuUserNumberOfPlusesVersion struct{ UInt32FieldVersion }
type PikabuUserNumberOfMinusesVersion struct{ UInt32FieldVersion }
type PikabuUserSignupTimestampVersion struct{ TimestampTypeFieldVersion }
type PikabuUserAvatarURLVersion struct{ StringFieldVersion }
type PikabuUserAwardsVersion struct {
	FieldVersionBase
	Value []uint64 `sql:",notnull,array"`
}
type PikabuUserApprovedTextVersion struct{ StringFieldVersion }
type PikabuUserCommunitiesVersion struct {
	FieldVersionBase
	Value []uint64 `sql:",notnull,array"`
}
type PikabuUserNumberOfSubscribersVersion struct{ UInt32FieldVersion }
type PikabuUserBanHistoryVersion struct {
	FieldVersionBase
	Value []uint64 `sql:",notnull,array"`
}
type PikabuUserBanEndTimestampVersion struct{ TimestampTypeFieldVersion }
type PikabuUserIsRatingHiddenVersion struct{ BoolFieldVersion }
type PikabuUserIsBannedVersion struct{ BoolFieldVersion }
type PikabuUserIsPermanentlyBannedVersion struct{ BoolFieldVersion }
type PikabuUserIsDeletedVersion struct{ BoolFieldVersion }

func init() {
	for _, item := range []interface{}{
		&PikabuUser{},
		&PikabuUserAward{},
		&PikabuUserCommunity{},
		&PikabuUserBanHistoryItem{},
		&PikabuUserUsernameVersion{},
		&PikabuUserRatingVersion{},
		&PikabuUserGenderVersion{},
		&PikabuUserNumberOfCommentsVersion{},
		&PikabuUserNumberOfStoriesVersion{},
		&PikabuUserNumberOfHotStoriesVersion{},
		&PikabuUserNumberOfPlusesVersion{},
		&PikabuUserNumberOfMinusesVersion{},
		&PikabuUserSignupTimestampVersion{},
		&PikabuUserAvatarURLVersion{},
		&PikabuUserAwardsVersion{},
		&PikabuUserApprovedTextVersion{},
		&PikabuUserCommunitiesVersion{},
		&PikabuUserNumberOfSubscribersVersion{},
		&PikabuUserBanHistoryVersion{},
		&PikabuUserBanEndTimestampVersion{},
		&PikabuUserIsRatingHiddenVersion{},
		&PikabuUserIsBannedVersion{},
		&PikabuUserIsPermanentlyBannedVersion{},
		&PikabuUserIsDeletedVersion{},
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
