package models

// generated code, do not touch!
// generated at timestamp 2020-01-01 12:13:25.223370339 &#43;0000 UTC m=&#43;0.002993701

type PikabuCommentParentIDVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentCreatedAtTimestampVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     TimestampType `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentTextVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentImagesVersion struct {
	ItemId    uint64               `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType        `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     []PikabuCommentImage `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentRatingVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentNumberOfPlusesVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentNumberOfMinusesVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentStoryIDVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentStoryURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentStoryTitleVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentAuthorIDVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentAuthorUsernameVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentAuthorGenderVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentAuthorAvatarURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentIgnoreCodeVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentIsIgnoredBySomeoneVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentIgnoredByVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     []string      `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentIsAuthorProfileDeletedVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentIsDeletedVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentIsAuthorCommunityModeratorVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentIsAuthorPikabuTeamVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentIsAuthorOfficialVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommentIsRatingHiddenVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityNameVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityLinkNameVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityAvatarURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityBackgroundImageURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityTagsVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     []string      `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityNumberOfStoriesVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityNumberOfSubscribersVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityDescriptionVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityRulesVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityRestrictionsVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityAdminIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuCommunityModeratorIdsVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     []uint64      `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryRatingVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryNumberOfPlusesVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryNumberOfMinusesVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryTitleVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryContentBlocksVersion struct {
	ItemId    uint64             `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType      `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     []PikabuStoryBlock `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryCreatedAtTimestampVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     TimestampType `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryStoryURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryTagsVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     []string      `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryNumberOfCommentsVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryIsDeletedVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryIsRatingHiddenVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryHasMineTagVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryHasAdultTagVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryIsLongpostVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryAuthorIDVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryAuthorUsernameVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryAuthorProfileURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryAuthorAvatarURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryCommunityLinkVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryCommunityNameVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryCommunityIDVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuStoryCommentsAreHotVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserUsernameVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserGenderVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserRatingVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserNumberOfCommentsVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserNumberOfSubscribersVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserNumberOfStoriesVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserNumberOfHotStoriesVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserNumberOfPlusesVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserNumberOfMinusesVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int32         `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserSignupTimestampVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     TimestampType `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAvatarURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserApprovedTextVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardIdsVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     []uint64      `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserCommunityIdsVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     []uint64      `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemIdsVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     []uint64      `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanEndTimestampVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     TimestampType `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserIsRatingHiddenVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserIsBannedVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserIsPermanentlyBannedVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardUserIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardAwardIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardAwardTitleVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardAwardImageURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardStoryIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardStoryTitleVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardIssuingDateVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardIsHiddenVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     bool          `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardCommentIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserAwardLinkVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemBanStartTimestampVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     TimestampType `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemCommentIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemCommentHtmlDeleteReasonVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemStoryIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemUserIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemBanReasonVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemBanReasonIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemStoryURLVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemModeratorIdVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemModeratorNameVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemModeratorAvatarVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemReasonsLimitVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemReasonCountVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     uint64        `sql:",notnull" json:"value" api:"order,filter"`
}

type PikabuUserBanHistoryItemReasonTitleVersion struct {
	ItemId    uint64        `sql:",pk,notnull" json:"item_id" api:"order,filter"`
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     string        `sql:",notnull" json:"value" api:"order,filter"`
}

var FieldsVersionTablesMap = map[string]interface{}{
	"PikabuCommentParentIDVersion":                           &PikabuCommentParentIDVersion{},
	"PikabuCommentCreatedAtTimestampVersion":                 &PikabuCommentCreatedAtTimestampVersion{},
	"PikabuCommentTextVersion":                               &PikabuCommentTextVersion{},
	"PikabuCommentImagesVersion":                             &PikabuCommentImagesVersion{},
	"PikabuCommentRatingVersion":                             &PikabuCommentRatingVersion{},
	"PikabuCommentNumberOfPlusesVersion":                     &PikabuCommentNumberOfPlusesVersion{},
	"PikabuCommentNumberOfMinusesVersion":                    &PikabuCommentNumberOfMinusesVersion{},
	"PikabuCommentStoryIDVersion":                            &PikabuCommentStoryIDVersion{},
	"PikabuCommentStoryURLVersion":                           &PikabuCommentStoryURLVersion{},
	"PikabuCommentStoryTitleVersion":                         &PikabuCommentStoryTitleVersion{},
	"PikabuCommentAuthorIDVersion":                           &PikabuCommentAuthorIDVersion{},
	"PikabuCommentAuthorUsernameVersion":                     &PikabuCommentAuthorUsernameVersion{},
	"PikabuCommentAuthorGenderVersion":                       &PikabuCommentAuthorGenderVersion{},
	"PikabuCommentAuthorAvatarURLVersion":                    &PikabuCommentAuthorAvatarURLVersion{},
	"PikabuCommentIgnoreCodeVersion":                         &PikabuCommentIgnoreCodeVersion{},
	"PikabuCommentIsIgnoredBySomeoneVersion":                 &PikabuCommentIsIgnoredBySomeoneVersion{},
	"PikabuCommentIgnoredByVersion":                          &PikabuCommentIgnoredByVersion{},
	"PikabuCommentIsAuthorProfileDeletedVersion":             &PikabuCommentIsAuthorProfileDeletedVersion{},
	"PikabuCommentIsDeletedVersion":                          &PikabuCommentIsDeletedVersion{},
	"PikabuCommentIsAuthorCommunityModeratorVersion":         &PikabuCommentIsAuthorCommunityModeratorVersion{},
	"PikabuCommentIsAuthorPikabuTeamVersion":                 &PikabuCommentIsAuthorPikabuTeamVersion{},
	"PikabuCommentIsAuthorOfficialVersion":                   &PikabuCommentIsAuthorOfficialVersion{},
	"PikabuCommentIsRatingHiddenVersion":                     &PikabuCommentIsRatingHiddenVersion{},
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
	"PikabuStoryRatingVersion":                               &PikabuStoryRatingVersion{},
	"PikabuStoryNumberOfPlusesVersion":                       &PikabuStoryNumberOfPlusesVersion{},
	"PikabuStoryNumberOfMinusesVersion":                      &PikabuStoryNumberOfMinusesVersion{},
	"PikabuStoryTitleVersion":                                &PikabuStoryTitleVersion{},
	"PikabuStoryContentBlocksVersion":                        &PikabuStoryContentBlocksVersion{},
	"PikabuStoryCreatedAtTimestampVersion":                   &PikabuStoryCreatedAtTimestampVersion{},
	"PikabuStoryStoryURLVersion":                             &PikabuStoryStoryURLVersion{},
	"PikabuStoryTagsVersion":                                 &PikabuStoryTagsVersion{},
	"PikabuStoryNumberOfCommentsVersion":                     &PikabuStoryNumberOfCommentsVersion{},
	"PikabuStoryIsDeletedVersion":                            &PikabuStoryIsDeletedVersion{},
	"PikabuStoryIsRatingHiddenVersion":                       &PikabuStoryIsRatingHiddenVersion{},
	"PikabuStoryHasMineTagVersion":                           &PikabuStoryHasMineTagVersion{},
	"PikabuStoryHasAdultTagVersion":                          &PikabuStoryHasAdultTagVersion{},
	"PikabuStoryIsLongpostVersion":                           &PikabuStoryIsLongpostVersion{},
	"PikabuStoryAuthorIDVersion":                             &PikabuStoryAuthorIDVersion{},
	"PikabuStoryAuthorUsernameVersion":                       &PikabuStoryAuthorUsernameVersion{},
	"PikabuStoryAuthorProfileURLVersion":                     &PikabuStoryAuthorProfileURLVersion{},
	"PikabuStoryAuthorAvatarURLVersion":                      &PikabuStoryAuthorAvatarURLVersion{},
	"PikabuStoryCommunityLinkVersion":                        &PikabuStoryCommunityLinkVersion{},
	"PikabuStoryCommunityNameVersion":                        &PikabuStoryCommunityNameVersion{},
	"PikabuStoryCommunityIDVersion":                          &PikabuStoryCommunityIDVersion{},
	"PikabuStoryCommentsAreHotVersion":                       &PikabuStoryCommentsAreHotVersion{},
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

var FieldsVersionAPITablesMap = map[string]interface{}{
	"PikabuCommentParentIDVersion":                           []PikabuCommentParentIDVersion{},
	"PikabuCommentCreatedAtTimestampVersion":                 []PikabuCommentCreatedAtTimestampVersion{},
	"PikabuCommentTextVersion":                               []PikabuCommentTextVersion{},
	"PikabuCommentImagesVersion":                             []PikabuCommentImagesVersion{},
	"PikabuCommentRatingVersion":                             []PikabuCommentRatingVersion{},
	"PikabuCommentNumberOfPlusesVersion":                     []PikabuCommentNumberOfPlusesVersion{},
	"PikabuCommentNumberOfMinusesVersion":                    []PikabuCommentNumberOfMinusesVersion{},
	"PikabuCommentStoryIDVersion":                            []PikabuCommentStoryIDVersion{},
	"PikabuCommentStoryURLVersion":                           []PikabuCommentStoryURLVersion{},
	"PikabuCommentStoryTitleVersion":                         []PikabuCommentStoryTitleVersion{},
	"PikabuCommentAuthorIDVersion":                           []PikabuCommentAuthorIDVersion{},
	"PikabuCommentAuthorUsernameVersion":                     []PikabuCommentAuthorUsernameVersion{},
	"PikabuCommentAuthorGenderVersion":                       []PikabuCommentAuthorGenderVersion{},
	"PikabuCommentAuthorAvatarURLVersion":                    []PikabuCommentAuthorAvatarURLVersion{},
	"PikabuCommentIgnoreCodeVersion":                         []PikabuCommentIgnoreCodeVersion{},
	"PikabuCommentIsIgnoredBySomeoneVersion":                 []PikabuCommentIsIgnoredBySomeoneVersion{},
	"PikabuCommentIgnoredByVersion":                          []PikabuCommentIgnoredByVersion{},
	"PikabuCommentIsAuthorProfileDeletedVersion":             []PikabuCommentIsAuthorProfileDeletedVersion{},
	"PikabuCommentIsDeletedVersion":                          []PikabuCommentIsDeletedVersion{},
	"PikabuCommentIsAuthorCommunityModeratorVersion":         []PikabuCommentIsAuthorCommunityModeratorVersion{},
	"PikabuCommentIsAuthorPikabuTeamVersion":                 []PikabuCommentIsAuthorPikabuTeamVersion{},
	"PikabuCommentIsAuthorOfficialVersion":                   []PikabuCommentIsAuthorOfficialVersion{},
	"PikabuCommentIsRatingHiddenVersion":                     []PikabuCommentIsRatingHiddenVersion{},
	"PikabuCommunityNameVersion":                             []PikabuCommunityNameVersion{},
	"PikabuCommunityLinkNameVersion":                         []PikabuCommunityLinkNameVersion{},
	"PikabuCommunityURLVersion":                              []PikabuCommunityURLVersion{},
	"PikabuCommunityAvatarURLVersion":                        []PikabuCommunityAvatarURLVersion{},
	"PikabuCommunityBackgroundImageURLVersion":               []PikabuCommunityBackgroundImageURLVersion{},
	"PikabuCommunityTagsVersion":                             []PikabuCommunityTagsVersion{},
	"PikabuCommunityNumberOfStoriesVersion":                  []PikabuCommunityNumberOfStoriesVersion{},
	"PikabuCommunityNumberOfSubscribersVersion":              []PikabuCommunityNumberOfSubscribersVersion{},
	"PikabuCommunityDescriptionVersion":                      []PikabuCommunityDescriptionVersion{},
	"PikabuCommunityRulesVersion":                            []PikabuCommunityRulesVersion{},
	"PikabuCommunityRestrictionsVersion":                     []PikabuCommunityRestrictionsVersion{},
	"PikabuCommunityAdminIdVersion":                          []PikabuCommunityAdminIdVersion{},
	"PikabuCommunityModeratorIdsVersion":                     []PikabuCommunityModeratorIdsVersion{},
	"PikabuStoryRatingVersion":                               []PikabuStoryRatingVersion{},
	"PikabuStoryNumberOfPlusesVersion":                       []PikabuStoryNumberOfPlusesVersion{},
	"PikabuStoryNumberOfMinusesVersion":                      []PikabuStoryNumberOfMinusesVersion{},
	"PikabuStoryTitleVersion":                                []PikabuStoryTitleVersion{},
	"PikabuStoryContentBlocksVersion":                        []PikabuStoryContentBlocksVersion{},
	"PikabuStoryCreatedAtTimestampVersion":                   []PikabuStoryCreatedAtTimestampVersion{},
	"PikabuStoryStoryURLVersion":                             []PikabuStoryStoryURLVersion{},
	"PikabuStoryTagsVersion":                                 []PikabuStoryTagsVersion{},
	"PikabuStoryNumberOfCommentsVersion":                     []PikabuStoryNumberOfCommentsVersion{},
	"PikabuStoryIsDeletedVersion":                            []PikabuStoryIsDeletedVersion{},
	"PikabuStoryIsRatingHiddenVersion":                       []PikabuStoryIsRatingHiddenVersion{},
	"PikabuStoryHasMineTagVersion":                           []PikabuStoryHasMineTagVersion{},
	"PikabuStoryHasAdultTagVersion":                          []PikabuStoryHasAdultTagVersion{},
	"PikabuStoryIsLongpostVersion":                           []PikabuStoryIsLongpostVersion{},
	"PikabuStoryAuthorIDVersion":                             []PikabuStoryAuthorIDVersion{},
	"PikabuStoryAuthorUsernameVersion":                       []PikabuStoryAuthorUsernameVersion{},
	"PikabuStoryAuthorProfileURLVersion":                     []PikabuStoryAuthorProfileURLVersion{},
	"PikabuStoryAuthorAvatarURLVersion":                      []PikabuStoryAuthorAvatarURLVersion{},
	"PikabuStoryCommunityLinkVersion":                        []PikabuStoryCommunityLinkVersion{},
	"PikabuStoryCommunityNameVersion":                        []PikabuStoryCommunityNameVersion{},
	"PikabuStoryCommunityIDVersion":                          []PikabuStoryCommunityIDVersion{},
	"PikabuStoryCommentsAreHotVersion":                       []PikabuStoryCommentsAreHotVersion{},
	"PikabuUserUsernameVersion":                              []PikabuUserUsernameVersion{},
	"PikabuUserGenderVersion":                                []PikabuUserGenderVersion{},
	"PikabuUserRatingVersion":                                []PikabuUserRatingVersion{},
	"PikabuUserNumberOfCommentsVersion":                      []PikabuUserNumberOfCommentsVersion{},
	"PikabuUserNumberOfSubscribersVersion":                   []PikabuUserNumberOfSubscribersVersion{},
	"PikabuUserNumberOfStoriesVersion":                       []PikabuUserNumberOfStoriesVersion{},
	"PikabuUserNumberOfHotStoriesVersion":                    []PikabuUserNumberOfHotStoriesVersion{},
	"PikabuUserNumberOfPlusesVersion":                        []PikabuUserNumberOfPlusesVersion{},
	"PikabuUserNumberOfMinusesVersion":                       []PikabuUserNumberOfMinusesVersion{},
	"PikabuUserSignupTimestampVersion":                       []PikabuUserSignupTimestampVersion{},
	"PikabuUserAvatarURLVersion":                             []PikabuUserAvatarURLVersion{},
	"PikabuUserApprovedTextVersion":                          []PikabuUserApprovedTextVersion{},
	"PikabuUserAwardIdsVersion":                              []PikabuUserAwardIdsVersion{},
	"PikabuUserCommunityIdsVersion":                          []PikabuUserCommunityIdsVersion{},
	"PikabuUserBanHistoryItemIdsVersion":                     []PikabuUserBanHistoryItemIdsVersion{},
	"PikabuUserBanEndTimestampVersion":                       []PikabuUserBanEndTimestampVersion{},
	"PikabuUserIsRatingHiddenVersion":                        []PikabuUserIsRatingHiddenVersion{},
	"PikabuUserIsBannedVersion":                              []PikabuUserIsBannedVersion{},
	"PikabuUserIsPermanentlyBannedVersion":                   []PikabuUserIsPermanentlyBannedVersion{},
	"PikabuUserAwardUserIdVersion":                           []PikabuUserAwardUserIdVersion{},
	"PikabuUserAwardAwardIdVersion":                          []PikabuUserAwardAwardIdVersion{},
	"PikabuUserAwardAwardTitleVersion":                       []PikabuUserAwardAwardTitleVersion{},
	"PikabuUserAwardAwardImageURLVersion":                    []PikabuUserAwardAwardImageURLVersion{},
	"PikabuUserAwardStoryIdVersion":                          []PikabuUserAwardStoryIdVersion{},
	"PikabuUserAwardStoryTitleVersion":                       []PikabuUserAwardStoryTitleVersion{},
	"PikabuUserAwardIssuingDateVersion":                      []PikabuUserAwardIssuingDateVersion{},
	"PikabuUserAwardIsHiddenVersion":                         []PikabuUserAwardIsHiddenVersion{},
	"PikabuUserAwardCommentIdVersion":                        []PikabuUserAwardCommentIdVersion{},
	"PikabuUserAwardLinkVersion":                             []PikabuUserAwardLinkVersion{},
	"PikabuUserBanHistoryItemBanStartTimestampVersion":       []PikabuUserBanHistoryItemBanStartTimestampVersion{},
	"PikabuUserBanHistoryItemCommentIdVersion":               []PikabuUserBanHistoryItemCommentIdVersion{},
	"PikabuUserBanHistoryItemCommentHtmlDeleteReasonVersion": []PikabuUserBanHistoryItemCommentHtmlDeleteReasonVersion{},
	"PikabuUserBanHistoryItemStoryIdVersion":                 []PikabuUserBanHistoryItemStoryIdVersion{},
	"PikabuUserBanHistoryItemUserIdVersion":                  []PikabuUserBanHistoryItemUserIdVersion{},
	"PikabuUserBanHistoryItemBanReasonVersion":               []PikabuUserBanHistoryItemBanReasonVersion{},
	"PikabuUserBanHistoryItemBanReasonIdVersion":             []PikabuUserBanHistoryItemBanReasonIdVersion{},
	"PikabuUserBanHistoryItemStoryURLVersion":                []PikabuUserBanHistoryItemStoryURLVersion{},
	"PikabuUserBanHistoryItemModeratorIdVersion":             []PikabuUserBanHistoryItemModeratorIdVersion{},
	"PikabuUserBanHistoryItemModeratorNameVersion":           []PikabuUserBanHistoryItemModeratorNameVersion{},
	"PikabuUserBanHistoryItemModeratorAvatarVersion":         []PikabuUserBanHistoryItemModeratorAvatarVersion{},
	"PikabuUserBanHistoryItemReasonsLimitVersion":            []PikabuUserBanHistoryItemReasonsLimitVersion{},
	"PikabuUserBanHistoryItemReasonCountVersion":             []PikabuUserBanHistoryItemReasonCountVersion{},
	"PikabuUserBanHistoryItemReasonTitleVersion":             []PikabuUserBanHistoryItemReasonTitleVersion{},
}

func init() {
	for _, item := range []interface{}{
		&PikabuCommentParentIDVersion{},
		&PikabuCommentCreatedAtTimestampVersion{},
		&PikabuCommentTextVersion{},
		&PikabuCommentImagesVersion{},
		&PikabuCommentRatingVersion{},
		&PikabuCommentNumberOfPlusesVersion{},
		&PikabuCommentNumberOfMinusesVersion{},
		&PikabuCommentStoryIDVersion{},
		&PikabuCommentStoryURLVersion{},
		&PikabuCommentStoryTitleVersion{},
		&PikabuCommentAuthorIDVersion{},
		&PikabuCommentAuthorUsernameVersion{},
		&PikabuCommentAuthorGenderVersion{},
		&PikabuCommentAuthorAvatarURLVersion{},
		&PikabuCommentIgnoreCodeVersion{},
		&PikabuCommentIsIgnoredBySomeoneVersion{},
		&PikabuCommentIgnoredByVersion{},
		&PikabuCommentIsAuthorProfileDeletedVersion{},
		&PikabuCommentIsDeletedVersion{},
		&PikabuCommentIsAuthorCommunityModeratorVersion{},
		&PikabuCommentIsAuthorPikabuTeamVersion{},
		&PikabuCommentIsAuthorOfficialVersion{},
		&PikabuCommentIsRatingHiddenVersion{},
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
		&PikabuStoryRatingVersion{},
		&PikabuStoryNumberOfPlusesVersion{},
		&PikabuStoryNumberOfMinusesVersion{},
		&PikabuStoryTitleVersion{},
		&PikabuStoryContentBlocksVersion{},
		&PikabuStoryCreatedAtTimestampVersion{},
		&PikabuStoryStoryURLVersion{},
		&PikabuStoryTagsVersion{},
		&PikabuStoryNumberOfCommentsVersion{},
		&PikabuStoryIsDeletedVersion{},
		&PikabuStoryIsRatingHiddenVersion{},
		&PikabuStoryHasMineTagVersion{},
		&PikabuStoryHasAdultTagVersion{},
		&PikabuStoryIsLongpostVersion{},
		&PikabuStoryAuthorIDVersion{},
		&PikabuStoryAuthorUsernameVersion{},
		&PikabuStoryAuthorProfileURLVersion{},
		&PikabuStoryAuthorAvatarURLVersion{},
		&PikabuStoryCommunityLinkVersion{},
		&PikabuStoryCommunityNameVersion{},
		&PikabuStoryCommunityIDVersion{},
		&PikabuStoryCommentsAreHotVersion{},
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
