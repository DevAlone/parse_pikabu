package models

// generated code, do not touch!
// generated at timestamp 2018-12-24 19:27:56.775710664 &#43;0000 UTC m=&#43;0.003865634

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

type PikabuUserAwardUserIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserAwardAwardIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserAwardAwardTitleVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserAwardAwardImageURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserAwardStoryIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserAwardStoryTitleVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserAwardIssuingDateVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserAwardIsHiddenVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     bool          `sql:",notnull"`
}

type PikabuUserAwardCommentIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserAwardLinkVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserBanHistoryItemBanStartTimestampVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     TimestampType `sql:",notnull"`
}

type PikabuUserBanHistoryItemCommentIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserBanHistoryItemCommentHtmlDeleteReasonVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserBanHistoryItemStoryIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserBanHistoryItemUserIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserBanHistoryItemBanReasonVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserBanHistoryItemBanReasonIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserBanHistoryItemStoryURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserBanHistoryItemModeratorIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserBanHistoryItemModeratorNameVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserBanHistoryItemModeratorAvatarVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

type PikabuUserBanHistoryItemReasonsLimitVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserBanHistoryItemReasonCountVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     uint64        `sql:",notnull"`
}

type PikabuUserBanHistoryItemReasonTitleVersion struct {
	ItemId    uint64        `sql:",pk,notnull"`
	Timestamp TimestampType `sql:",pk,notnull"`
	Value     string        `sql:",notnull"`
}

var FieldsVersionTablesMap = map[string]interface{}{
	"PikabuCommunityNameVersion":                             &PikabuCommunityNameVersion{},
	"PikabuCommunityLinkNameVersion":                         &PikabuCommunityLinkNameVersion{},
	"PikabuCommunityURLVersion":                              &PikabuCommunityURLVersion{},
	"PikabuCommunityAvatarURLVersion":                        &PikabuCommunityAvatarURLVersion{},
	"PikabuCommunityBackgroundImageURLVersion":               &PikabuCommunityBackgroundImageURLVersion{},
	"PikabuCommunityTagsVersion":                             &PikabuCommunityTagsVersion{},
	"PikabuCommunityNumberOfStoriesVersion":                  &PikabuCommunityNumberOfStoriesVersion{},
	"PikabuCommunityNumberOfSubscribersVersion":              &PikabuCommunityNumberOfSubscribersVersion{},
	"PikabuCommunityDescriptionVersion":                      &PikabuCommunityDescriptionVersion{},
	"PikabuCommunityRulesVersion":                            &PikabuCommunityRulesVersion{},
	"PikabuCommunityRestrictionsVersion":                     &PikabuCommunityRestrictionsVersion{},
	"PikabuCommunityAdminIdVersion":                          &PikabuCommunityAdminIdVersion{},
	"PikabuCommunityModeratorIdsVersion":                     &PikabuCommunityModeratorIdsVersion{},
	"PikabuUserUsernameVersion":                              &PikabuUserUsernameVersion{},
	"PikabuUserGenderVersion":                                &PikabuUserGenderVersion{},
	"PikabuUserRatingVersion":                                &PikabuUserRatingVersion{},
	"PikabuUserNumberOfCommentsVersion":                      &PikabuUserNumberOfCommentsVersion{},
	"PikabuUserNumberOfSubscribersVersion":                   &PikabuUserNumberOfSubscribersVersion{},
	"PikabuUserNumberOfStoriesVersion":                       &PikabuUserNumberOfStoriesVersion{},
	"PikabuUserNumberOfHotStoriesVersion":                    &PikabuUserNumberOfHotStoriesVersion{},
	"PikabuUserNumberOfPlusesVersion":                        &PikabuUserNumberOfPlusesVersion{},
	"PikabuUserNumberOfMinusesVersion":                       &PikabuUserNumberOfMinusesVersion{},
	"PikabuUserSignupTimestampVersion":                       &PikabuUserSignupTimestampVersion{},
	"PikabuUserAvatarURLVersion":                             &PikabuUserAvatarURLVersion{},
	"PikabuUserApprovedTextVersion":                          &PikabuUserApprovedTextVersion{},
	"PikabuUserAwardIdsVersion":                              &PikabuUserAwardIdsVersion{},
	"PikabuUserCommunityIdsVersion":                          &PikabuUserCommunityIdsVersion{},
	"PikabuUserBanHistoryItemIdsVersion":                     &PikabuUserBanHistoryItemIdsVersion{},
	"PikabuUserBanEndTimestampVersion":                       &PikabuUserBanEndTimestampVersion{},
	"PikabuUserIsRatingHiddenVersion":                        &PikabuUserIsRatingHiddenVersion{},
	"PikabuUserIsBannedVersion":                              &PikabuUserIsBannedVersion{},
	"PikabuUserIsPermanentlyBannedVersion":                   &PikabuUserIsPermanentlyBannedVersion{},
	"PikabuUserAwardUserIdVersion":                           &PikabuUserAwardUserIdVersion{},
	"PikabuUserAwardAwardIdVersion":                          &PikabuUserAwardAwardIdVersion{},
	"PikabuUserAwardAwardTitleVersion":                       &PikabuUserAwardAwardTitleVersion{},
	"PikabuUserAwardAwardImageURLVersion":                    &PikabuUserAwardAwardImageURLVersion{},
	"PikabuUserAwardStoryIdVersion":                          &PikabuUserAwardStoryIdVersion{},
	"PikabuUserAwardStoryTitleVersion":                       &PikabuUserAwardStoryTitleVersion{},
	"PikabuUserAwardIssuingDateVersion":                      &PikabuUserAwardIssuingDateVersion{},
	"PikabuUserAwardIsHiddenVersion":                         &PikabuUserAwardIsHiddenVersion{},
	"PikabuUserAwardCommentIdVersion":                        &PikabuUserAwardCommentIdVersion{},
	"PikabuUserAwardLinkVersion":                             &PikabuUserAwardLinkVersion{},
	"PikabuUserBanHistoryItemBanStartTimestampVersion":       &PikabuUserBanHistoryItemBanStartTimestampVersion{},
	"PikabuUserBanHistoryItemCommentIdVersion":               &PikabuUserBanHistoryItemCommentIdVersion{},
	"PikabuUserBanHistoryItemCommentHtmlDeleteReasonVersion": &PikabuUserBanHistoryItemCommentHtmlDeleteReasonVersion{},
	"PikabuUserBanHistoryItemStoryIdVersion":                 &PikabuUserBanHistoryItemStoryIdVersion{},
	"PikabuUserBanHistoryItemUserIdVersion":                  &PikabuUserBanHistoryItemUserIdVersion{},
	"PikabuUserBanHistoryItemBanReasonVersion":               &PikabuUserBanHistoryItemBanReasonVersion{},
	"PikabuUserBanHistoryItemBanReasonIdVersion":             &PikabuUserBanHistoryItemBanReasonIdVersion{},
	"PikabuUserBanHistoryItemStoryURLVersion":                &PikabuUserBanHistoryItemStoryURLVersion{},
	"PikabuUserBanHistoryItemModeratorIdVersion":             &PikabuUserBanHistoryItemModeratorIdVersion{},
	"PikabuUserBanHistoryItemModeratorNameVersion":           &PikabuUserBanHistoryItemModeratorNameVersion{},
	"PikabuUserBanHistoryItemModeratorAvatarVersion":         &PikabuUserBanHistoryItemModeratorAvatarVersion{},
	"PikabuUserBanHistoryItemReasonsLimitVersion":            &PikabuUserBanHistoryItemReasonsLimitVersion{},
	"PikabuUserBanHistoryItemReasonCountVersion":             &PikabuUserBanHistoryItemReasonCountVersion{},
	"PikabuUserBanHistoryItemReasonTitleVersion":             &PikabuUserBanHistoryItemReasonTitleVersion{},
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
		&PikabuUserAwardUserIdVersion{},
		&PikabuUserAwardAwardIdVersion{},
		&PikabuUserAwardAwardTitleVersion{},
		&PikabuUserAwardAwardImageURLVersion{},
		&PikabuUserAwardStoryIdVersion{},
		&PikabuUserAwardStoryTitleVersion{},
		&PikabuUserAwardIssuingDateVersion{},
		&PikabuUserAwardIsHiddenVersion{},
		&PikabuUserAwardCommentIdVersion{},
		&PikabuUserAwardLinkVersion{},
		&PikabuUserBanHistoryItemBanStartTimestampVersion{},
		&PikabuUserBanHistoryItemCommentIdVersion{},
		&PikabuUserBanHistoryItemCommentHtmlDeleteReasonVersion{},
		&PikabuUserBanHistoryItemStoryIdVersion{},
		&PikabuUserBanHistoryItemUserIdVersion{},
		&PikabuUserBanHistoryItemBanReasonVersion{},
		&PikabuUserBanHistoryItemBanReasonIdVersion{},
		&PikabuUserBanHistoryItemStoryURLVersion{},
		&PikabuUserBanHistoryItemModeratorIdVersion{},
		&PikabuUserBanHistoryItemModeratorNameVersion{},
		&PikabuUserBanHistoryItemModeratorAvatarVersion{},
		&PikabuUserBanHistoryItemReasonsLimitVersion{},
		&PikabuUserBanHistoryItemReasonCountVersion{},
		&PikabuUserBanHistoryItemReasonTitleVersion{},
	} {
		Tables = append(Tables, item)
	}
}
