package models

// PikabuCommentImage // TODO: add doc
type PikabuCommentImage struct {
	SmallURL            string            `json:"small_url"`
	LargeURL            string            `json:"large_url"`
	AnimationPreviewURL string            `json:"animation_preview_url"`
	AnimationBaseURL    string            `json:"animation_base_url"`
	AnimationFormats    map[string]uint64 `json:"animation_formats"`
	Size                []uint64          `json:"img_size"`
}

// PikabuComment // TODO: add doc
type PikabuComment struct {
	PikabuID uint64 `sql:",pk" json:"pikabu_id" api:"order,filter"`

	ParentID           uint64        `sql:",notnull" json:"parent_id" gen_versions:""`
	CreatedAtTimestamp TimestampType `sql:",notnull" json:"created_at_timestamp" gen_versions:"" gen_distributions:"86400"`

	Text                   string               `sql:",notnull" json:"text" gen_versions:""`
	Images                 []PikabuCommentImage `sql:",notnull" json:"images" gen_versions:""`
	Rating                 int32                `sql:",notnull" json:"rating" gen_versions:""`
	NumberOfPluses         int32                `sql:",notnull" json:"number_of_pluses" gen_versions:""`
	NumberOfMinuses        int32                `sql:",notnull" json:"number_of_minuses" gen_versions:""`
	StoryID                uint64               `sql:",notnull" json:"story_id" gen_versions:""`
	StoryURL               string               `sql:",notnull" json:"story_url" gen_versions:""`
	StoryTitle             string               `sql:",notnull" json:"story_title" gen_versions:""`
	AuthorID               uint64               `sql:",notnull" json:"author_id" gen_versions:""`
	AuthorUsername         string               `sql:",notnull" json:"author_username" gen_versions:""`
	AuthorGender           int32                `sql:",notnull" json:"author_gender" gen_versions:""`
	AuthorAvatarURL        string               `sql:",notnull" json:"author_avatar_url" gen_versions:""`
	IgnoreCode             int32                `sql:",notnull" json:"ignore_code" gen_versions:""`
	IsIgnoredBySomeone     bool                 `sql:",notnull" json:"is_ignored_by_someone" gen_versions:""`
	IgnoredBy              []string             `sql:",notnull" json:"ignored_by" gen_versions:""`
	IsAuthorProfileDeleted bool                 `sql:",notnull" json:"is_author_profile_deleted" gen_versions:""`
	IsDeleted              bool                 `sql:",notnull" json:"is_deleted" gen_versions:""`

	IsAuthorCommunityModerator bool `sql:",notnull" json:"is_author_community_moderator" gen_versions:""`
	IsAuthorPikabuTeam         bool `sql:",notnull" json:"is_author_pikabu_team" gen_versions:""`
	IsAuthorOfficial           bool `sql:",notnull" json:"is_author_official" gen_versions:""`
	IsRatingHidden             bool `sql:",notnull" json:"is_rating_hidden" gen_versions:""`

	AddedTimestamp       TimestampType `sql:",notnull" json:"added_timestamp" api:"order"`
	LastUpdateTimestamp  TimestampType `sql:",notnull" json:"last_update_timestamp" api:"order"`
	NextUpdateTimestamp  TimestampType `sql:",notnull" json:"next_update_timestamp" api:"order"`
	TaskTakenAtTimestamp TimestampType `sql:",notnull" json:"task_taken_at_timestamp" api:"order"`
}

// PikabuDeletedOrNeverExistedComment // TODO: add doc
type PikabuDeletedOrNeverExistedComment struct {
	PikabuID             uint64        `sql:",pk" json:"pikabu_id"`
	LastUpdateTimestamp  TimestampType `sql:",notnull" json:"last_update_timestamp"`
	NextUpdateTimestamp  TimestampType `sql:",notnull" json:"next_update_timestamp"`
	TaskTakenAtTimestamp TimestampType `sql:",notnull" json:"task_taken_at_timestamp" api:"order"`
}

func init() {
	Tables = append(Tables, []interface{}{
		&PikabuComment{},
		&PikabuDeletedOrNeverExistedComment{},
	}...)

	addIndex("pikabu_comments", "parent_id", "")
	// TODO: addIndex("pikabu_comments", "parent_id", "hash") ?
	addIndex("pikabu_comments", "created_at_timestamp", "")
	addIndex("pikabu_comments", "text gin_trgm_ops", "gin")
	// TODO: Text                   ?
	// TODO: Images                ?
	addIndex("pikabu_comments", "rating", "")
	addIndex("pikabu_comments", "number_of_pluses", "")
	addIndex("pikabu_comments", "number_of_minuses", "")
	addIndex("pikabu_comments", "story_id", "")
	// TODO: StoryTitle      ?
	addIndex("pikabu_comments", "author_id", "")
	addIndex("pikabu_comments", "author_username", "")
	// TODO: AuthorGender  ?
	addIndex("pikabu_comments", "ignore_code", "")
	addIndex("pikabu_comments", "is_ignored_by_someone", "")
	// TODO: IgnoredBy ?
	addIndex("pikabu_comments", "is_author_profile_deleted", "")
	addIndex("pikabu_comments", "is_deleted", "")
	addIndex("pikabu_comments", "is_author_community_moderator", "")
	addIndex("pikabu_comments", "is_author_pikabu_team", "")
	addIndex("pikabu_comments", "is_author_official", "")
	addIndex("pikabu_comments", "is_rating_hidden", "")

	addIndex("pikabu_comments", "added_timestamp", "")
	addIndex("pikabu_comments", "last_update_timestamp", "")
	addIndex("pikabu_comments", "next_update_timestamp", "")
	addIndex("pikabu_comments", "task_taken_at_timestamp", "")

	addIndex("pikabu_deleted_or_never_existed_comments", "last_update_timestamp", "")
	addIndex("pikabu_deleted_or_never_existed_comments", "next_update_timestamp", "")
	addIndex("pikabu_deleted_or_never_existed_comments", "task_taken_at_timestamp", "")
}
