package models

// TODO: update structure

type User struct {
	PikabuId uint64 `sql:",pk"`

	Username           string        `sql:",notnull"`
	Rating             int32         `sql:",notnull"`
	Gender             string        `sql:",notnull"`
	NumberOfComments   uint32        `sql:",notnull"`
	NumberOfStories    uint32        `sql:",notnull"`
	NumberOfHotStories uint32        `sql:",notnull"`
	NumberOfPluses     uint32        `sql:",notnull"`
	NumberOfMinuses    uint32        `sql:",notnull"`
	SignupTimestamp    TimestampType `sql:",notnull"`
	AvatarURL          string        `sql:",notnull"`
	// Awards             string
	Awards
	ApprovedText string `sql:",notnull"`
	Communities
	NumberOfSubscribers uint32 `sql:",notnull"`
	BanHistory
	BanEndTimestamp     TimestampType `sql:",notnull"`
	IsRatingHidden      bool          `sql:",notnull"`
	IsBanned            bool          `sql:",notnull"`
	IsPermanentlyBanned bool          `sql:",notnull"`

	// ?
	IsDeleted bool `sql:",notnull,default:false"`

	LastUpdateTimestamp TimestampType `sql:",notnull"`
	NextUpdateTimestamp TimestampType `sql:",notnull"`
}

type UserUsernameVersion struct{ StringFieldVersion}
type UserRatingVersion struct { Int32FieldVersion }
type UserGenderVersion struct {StringFieldVersion}
type UserNumberOfCommentsVersion struct { UInt32FieldVersion }
type UserNumberOfStoriesVersion struct {UInt32FieldVersion}
type UserNumberOfHotStoriesVersion struct {UInt32FieldVersion}
type UserNumberOfPlusesVersion struct {UInt32FieldVersion }
type UserNumberOfMinusesVersion    struct { UInt32FieldVersion }
type UserSignupTimestampVersion    struct { TimestampTypeFieldVersion}
type UserAvatarURLVersion struct {StringFieldVersion }
Awards
type UserApprovedTextVersion struct { StringFieldVersion }
Communities
type UserNumberOfSubscribersVersion struct {UInt32FieldVersion }
BanHistory
type UserBanEndTimestampVersion struct { TimestampTypeFieldVersion}
type UserIsRatingHiddenVersion struct {BoolFieldVersion }
type UserIsBannedVersion struct {BoolFieldVersion }
type UserIsPermanentlyBannedVersion struct {BoolFieldVersion }
type UserIsDeletedVersion struct {BoolFieldVersion}

func init() {
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
