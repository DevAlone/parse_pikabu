package models

// PikabuStoryBlock // TODO: add doc
type PikabuStoryBlock struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// PikabuStory // TODO: add doc
type PikabuStory struct {
	PikabuID uint64 `sql:",pk" json:"pikabu_id" api:"order,filter"`

	Rating             int32              `sql:",notnull" json:"rating" gen_versions:"" api:"order,filter"`
	NumberOfPluses     int32              `sql:",notnull" json:"number_of_pluses" gen_versions:"" api:"order,filter"`
	NumberOfMinuses    int32              `sql:",notnull" json:"number_of_minuses" gen_versions:"" api:"order,filter"`
	Title              string             `sql:",notnull" json:"title" gen_versions:"" api:"order,filter"`
	ContentBlocks      []PikabuStoryBlock `sql:",notnull" json:"content_blocks" gen_versions:""`
	CreatedAtTimestamp TimestampType      `sql:",notnull" json:"created_at_timestamp" gen_versions:"" gen_distributions:"86400" api:"order,filter"`
	StoryURL           string             `sql:",notnull" json:"story_url" gen_versions:""`
	Tags               []string           `sql:",notnull" json:"tags" gen_versions:""`
	NumberOfComments   int32              `sql:",notnull" json:"number_of_comments" gen_versions:"" api:"order,filter"`
	IsDeleted          bool               `sql:",notnull" json:"is_deleted" gen_versions:"" api:"order,filter"`
	IsRatingHidden     bool               `sql:",notnull" json:"is_rating_hidden" gen_versions:"" api:"order,filter"`
	HasMineTag         bool               `sql:",notnull" json:"has_mine_tag" gen_versions:"" api:"order,filter"`
	HasAdultTag        bool               `sql:",notnull" json:"has_adult_tag" gen_versions:"" api:"order,filter"`
	IsLongpost         bool               `sql:",notnull" json:"is_longpost" gen_versions:"" api:"order,filter"`
	AuthorID           uint64             `sql:",notnull" json:"author_id" gen_versions:"" api:"order,filter"`
	AuthorUsername     string             `sql:",notnull" json:"author_username" gen_versions:"" api:"order,filter"`
	AuthorProfileURL   string             `sql:",notnull" json:"author_profile_url" gen_versions:""`
	AuthorAvatarURL    string             `sql:",notnull" json:"author_avatar_url" gen_versions:""`
	CommunityLink      string             `sql:",notnull" json:"community_link" gen_versions:"" api:"order,filter"`
	CommunityName      string             `sql:",notnull" json:"community_name" gen_versions:""`
	CommunityID        uint64             `sql:",notnull" json:"community_id" gen_versions:"" api:"order,filter"`
	CommentsAreHot     bool               `sql:",notnull" json:"comments_are_hot" gen_versions:"" api:"order,filter"`

	AddedTimestamp       TimestampType `sql:",notnull" json:"added_timestamp" api:"order,filter"`
	LastUpdateTimestamp  TimestampType `sql:",notnull" json:"last_update_timestamp" api:"order,filter"`
	NextUpdateTimestamp  TimestampType `sql:",notnull" json:"next_update_timestamp" api:"order,filter"`
	TaskTakenAtTimestamp TimestampType `sql:",notnull" json:"task_taken_at_timestamp" api:"order,filter"`
	IsPermanentlyDeleted bool          `sql:",notnull,default:false" json:"is_permanently_deleted" api:"order,filter"`
	IsHiddenInAPI        bool          `sql:",notnull,default:false" json:"is_hidden_in_api" api:"order,filter"`
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

	addIndex("pikabu_stories", "title", "")
	addIndex("pikabu_stories", "LOWER(title)", "hash")
	addIndex("pikabu_stories", "title gin_trgm_ops", "gin")

	addIndex("pikabu_stories", "created_at_timestamp", "")
	addIndex("pikabu_stories", "number_of_comments", "")
	addIndex("pikabu_stories", "is_deleted", "")
	addIndex("pikabu_stories", "is_rating_hidden", "")
	addIndex("pikabu_stories", "has_mine_tag", "")
	addIndex("pikabu_stories", "has_adult_tag", "")
	addIndex("pikabu_stories", "is_longpost", "")
	addIndex("pikabu_stories", "author_id", "")

	addIndex("pikabu_stories", "author_username", "")
	addIndex("pikabu_stories", "LOWER(author_username)", "hash")

	addIndex("pikabu_stories", "community_link", "")
	addIndex("pikabu_stories", "LOWER(community_link)", "hash")

	// TODO: community_name?

	addIndex("pikabu_stories", "community_id", "")
	addIndex("pikabu_stories", "comments_are_hot", "")
	addIndex("pikabu_stories", "added_timestamp", "")
	addIndex("pikabu_stories", "last_update_timestamp", "")
	addIndex("pikabu_stories", "next_update_timestamp", "")
	addIndex("pikabu_stories", "task_taken_at_timestamp", "")
	addIndex("pikabu_stories", "is_permanently_deleted", "")
	addIndex("pikabu_stories", "is_hidden_in_api", "")

	addIndex("pikabu_deleted_or_never_existed_stories", "last_update_timestamp", "")
	addIndex("pikabu_deleted_or_never_existed_stories", "next_update_timestamp", "")
	addIndex("pikabu_deleted_or_never_existed_stories", "task_taken_at_timestamp", "")
}
