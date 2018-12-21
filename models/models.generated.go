package models

// generated code, do not touch!
// generated at timestamp 2018-12-21 19:51:32.629609812 &#43;0000 UTC m=&#43;0.003748290

type PikabuCommunityNameVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuCommunityLinkNameVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuCommunityURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuCommunityAvatarURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuCommunityBackgroundImageURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuCommunityTagsVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     []string      `sql:",notnull"`
}

type PikabuCommunityNumberOfStoriesVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     int32         `sql:",notnull"`
}

type PikabuCommunityNumberOfSubscribersVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     int32         `sql:",notnull"`
}

type PikabuCommunityDescriptionVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuCommunityRulesVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuCommunityRestrictionsVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuCommunityAdminIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuCommunityModeratorIdsVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     []uint64      `sql:",notnull"`
}

type PikabuUserUsernameVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserGenderVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserRatingVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     float32       `sql:",notnull"`
}

type PikabuUserNumberOfCommentsVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     int32         `sql:",notnull"`
}

type PikabuUserNumberOfSubscribersVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     int32         `sql:",notnull"`
}

type PikabuUserNumberOfStoriesVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     int32         `sql:",notnull"`
}

type PikabuUserNumberOfHotStoriesVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     int32         `sql:",notnull"`
}

type PikabuUserNumberOfPlusesVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     int32         `sql:",notnull"`
}

type PikabuUserNumberOfMinusesVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     int32         `sql:",notnull"`
}

type PikabuUserSignupTimestampVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     TimestampType `sql:",notnull"`
}

type PikabuUserAvatarURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserApprovedTextVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserAwardIdsVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     []uint64      `sql:",notnull"`
}

type PikabuUserCommunityIdsVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     []uint64      `sql:",notnull"`
}

type PikabuUserBanHistoryItemIdsVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     []uint64      `sql:",notnull"`
}

type PikabuUserBanEndTimestampVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     TimestampType `sql:",notnull"`
}

type PikabuUserIsRatingHiddenVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     bool          `sql:",notnull"`
}

type PikabuUserIsBannedVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     bool          `sql:",notnull"`
}

type PikabuUserIsPermanentlyBannedVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     bool          `sql:",notnull"`
}

func init() {
	for _, item := range []interface{}{
		&PikabuCommunityNameVersion{},
		&PikabuCommunityLinkNameVersion{},
		&PikabuCommunityURLVersion{},
		&PikabuCommunityAvatarURLVersion{},
		&PikabuCommunityBackgroundImageURLVersion{},
		&PikabuCommunityTagsVersion{},
		&PikabuCommunityNumberOfStoriesVersion{},
		&PikabuCommunityNumberOfSubscribersVersion{},
		&PikabuCommunityDescriptionVersion{},
		&PikabuCommunityRulesVersion{},
		&PikabuCommunityRestrictionsVersion{},
		&PikabuCommunityAdminIdVersion{},
		&PikabuCommunityModeratorIdsVersion{},
		&PikabuUserUsernameVersion{},
		&PikabuUserGenderVersion{},
		&PikabuUserRatingVersion{},
		&PikabuUserNumberOfCommentsVersion{},
		&PikabuUserNumberOfSubscribersVersion{},
		&PikabuUserNumberOfStoriesVersion{},
		&PikabuUserNumberOfHotStoriesVersion{},
		&PikabuUserNumberOfPlusesVersion{},
		&PikabuUserNumberOfMinusesVersion{},
		&PikabuUserSignupTimestampVersion{},
		&PikabuUserAvatarURLVersion{},
		&PikabuUserApprovedTextVersion{},
		&PikabuUserAwardIdsVersion{},
		&PikabuUserCommunityIdsVersion{},
		&PikabuUserBanHistoryItemIdsVersion{},
		&PikabuUserBanEndTimestampVersion{},
		&PikabuUserIsRatingHiddenVersion{},
		&PikabuUserIsBannedVersion{},
		&PikabuUserIsPermanentlyBannedVersion{},
	} {
		Tables = append(Tables, item)
	}
}
