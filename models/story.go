package models

// PikabuStoryBlock // TODO: add doc
type PikabuStoryBlock struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// PikabuStory // TODO: add doc
type PikabuStory struct {
	PikabuID uint64 `sql:",pk" json:"pikabu_id" api:"order,filter"`

	Rating             int32              `sql:",notnull" json:"rating" gen_versions:""`
	NumberOfPluses     int32              `sql:",notnull" json:"number_of_pluses" gen_versions:""`
	NumberOfMinuses    int32              `sql:",notnull" json:"number_of_minuses" gen_versions:""`
	Title              string             `sql:",notnull" json:"title" gen_versions:""`
	ContentBlocks      []PikabuStoryBlock `sql:",notnull" json:"content_blocks" gen_versions:""`
	CreatedAtTimestamp TimestampType      `sql:",notnull" json:"created_at_timestamp" gen_versions:""`
	StoryURL           string             `sql:",notnull" json:"story_url" gen_versions:""`
	Tags               []string           `sql:",notnull" json:"tags" gen_versions:""`
	NumberOfComments   int32              `sql:",notnull" json:"number_of_comments" gen_versions:""`
	IsDeleted          bool               `sql:",notnull" json:"is_deleted" gen_versions:""`
	IsRatingHidden     bool               `sql:",notnull" json:"is_rating_hidden" gen_versions:""`
	HasMineTag         bool               `sql:",notnull" json:"has_mine_tag" gen_versions:""`
	HasAdultTag        bool               `sql:",notnull" json:"has_adult_tag" gen_versions:""`
	IsLongpost         bool               `sql:",notnull" json:"is_longpost" gen_versions:""`
	AuthorID           uint64             `sql:",notnull" json:"author_id" gen_versions:""`
	AuthorUsername     string             `sql:",notnull" json:"author_username" gen_versions:""`
	AuthorProfileURL   string             `sql:",notnull" json:"author_profile_url" gen_versions:""`
	AuthorAvatarURL    string             `sql:",notnull" json:"author_avatar_url" gen_versions:""`
	CommunityLink      string             `sql:",notnull" json:"community_link" gen_versions:""`
	CommunityName      string             `sql:",notnull" json:"community_name" gen_versions:""`
	CommunityID        uint64             `sql:",notnull" json:"community_id" gen_versions:""`
	CommentsAreHot     bool               `sql:",notnull" json:"comments_are_hot" gen_versions:""`

	AddedTimestamp       TimestampType `sql:",notnull" json:"added_timestamp" api:"order"`
	LastUpdateTimestamp  TimestampType `sql:",notnull" json:"last_update_timestamp" api:"order"`
	NextUpdateTimestamp  TimestampType `sql:",notnull" json:"next_update_timestamp" api:"order"`
	TaskTakenAtTimestamp TimestampType `sql:",notnull" json:"task_taken_at_timestamp" api:"order"`
	IsPermanentlyDeleted bool          `sql:",notnull,default:false" json:"is_permanently_deleted"`
}

// PikabuDeletedOrNeverExistedStory // TODO: add doc
type PikabuDeletedOrNeverExistedStory struct {
	PikabuID             uint64        `sql:",pk" json:"pikabu_id"`
	LastUpdateTimestamp  TimestampType `sql:",notnull" json:"last_update_timestamp"`
	NextUpdateTimestamp  TimestampType `sql:",notnull" json:"next_update_timestamp"`
	TaskTakenAtTimestamp TimestampType `sql:",notnull" json:"task_taken_at_timestamp" api:"order"`
}

func init() {
	Tables = append(Tables, []interface{}{
		&PikabuStory{},
		&PikabuDeletedOrNeverExistedStory{},
	}...)

	addIndex("pikabu_stories", "rating", "")
	addIndex("pikabu_stories", "number_of_pluses", "")
	addIndex("pikabu_stories", "number_of_minuses", "")
	// TODO: addIndex("pikabu_stories", "title", "") ?
	addIndex("pikabu_stories", "created_at_timestamp", "")
	addIndex("pikabu_stories", "number_of_comments", "")
	addIndex("pikabu_stories", "is_deleted", "")
	addIndex("pikabu_stories", "is_rating_hidden", "")
	addIndex("pikabu_stories", "has_mine_tag", "")
	addIndex("pikabu_stories", "has_adult_tag", "")
	addIndex("pikabu_stories", "is_longpost", "")
	addIndex("pikabu_stories", "author_id", "")
	// addIndex("pikabu_stories", "author_username", "hash")
	addIndex("pikabu_stories", "LOWER(author_username)", "hash")
	addIndex("pikabu_stories", "LOWER(community_link)", "hash")
	// addIndex("pikabu_stories", "community_name", "")
	addIndex("pikabu_stories", "comments_are_hot", "")

	addIndex("pikabu_stories", "added_timestamp", "")
	addIndex("pikabu_stories", "last_update_timestamp", "")
	addIndex("pikabu_stories", "next_update_timestamp", "")
	addIndex("pikabu_stories", "task_taken_at_timestamp", "")
	addIndex("pikabu_stories", "is_permanently_deleted", "")

	addIndex("pikabu_deleted_or_never_existed_stories", "last_update_timestamp", "")
	addIndex("pikabu_deleted_or_never_existed_stories", "next_update_timestamp", "")
	addIndex("pikabu_deleted_or_never_existed_stories", "task_taken_at_timestamp", "")
}
