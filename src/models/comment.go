package models

// all fields are the last versions of them
type Comment struct {
	Id       uint64
	ParentId uint64
	// Parent                     *Comment // TODO: check why this doesn't work // ERROR #23503 insert or update on table "comments" violates foreign key constraint "comments_parent_id_fkey
	Children              []Comment `pg:"fk:parent_id"`
	CreationTimestamp     int32 `sql:",notnull"`
	FirstParsingTimestamp int32 `sql:",notnull"`
	LastParsingTimestamp  int32 `sql:",notnull"`
	// can be null
	Rating                int32 `sql:",notnull"`
	StoryId               uint64 `sql:",notnull"` // TODO: add foreign key for story
	UserId                int32 `sql:",notnull"`
	// User                       *User
	AuthorUsername             string `sql:",notnull"`
	IsHidden                   bool `sql:",notnull"`
	IsDeleted                  bool `sql:",notnull"`
	IsAuthorCommunityModerator bool `sql:",notnull"`
	IsAuthorPikabuTeam         bool `sql:",notnull"`

	// TODO: Content
	Text string `sql:",notnull"`
}

type CommentImagesVersion struct {
	ParsingTimestamp int32  `sql:",pk,notnull"`
	CommentId        uint64 `sql:",pk,notnull"`
	Comment          *Comment
	ImageIds []uint64 `pg:",array"`
}

type CommentParentIdVersion struct{ Uint64FieldVersion }
type CommentCreatingTimestampVersion struct{ Int32FieldVersion }
type CommentRatingVersion struct{ Int32FieldVersion }
type CommentStoryIdVersion struct{ Uint64FieldVersion }
type CommentUserIdVersion struct{ Int32FieldVersion }
type CommentAuthorUsernameVersion struct{ TextFieldVersion }
type CommentIsHiddenVersion struct{ BoolFieldVersion }
type CommentIsDeletedVersion struct{ BoolFieldVersion }
type CommentIsAuthorCommunityModeratorVersion struct{ BoolFieldVersion }
type CommentIsAuthorPikabuTeamVersion struct{ BoolFieldVersion }

type CommentTextVersion struct {
	CommentId uint64 `sql:",pk,notnull"`
	Timestamp int32  `sql:",pk,notnull"`
	Diffs     string `sql:",notnull"`
}

func init() {
	addIndex("comments", "parent_id", "")
	addIndex("comments", "creation_timestamp", "")
	addIndex("comments", "first_parsing_timestamp", "")
	addIndex("comments", "last_parsing_timestamp", "")
	addIndex("comments", "rating", "")
	addIndex("comments", "story_id", "")
	addIndex("comments", "user_id", "")
	addIndex("comments", "author_username", "hash")
	addIndex("comments", "is_hidden", "")
	addIndex("comments", "is_deleted", "")
	addIndex("comments", "is_author_community_moderator", "")
	addIndex("comments", "is_author_pikabu_team", "")

	addIndex("comment_images_versions", "comment_id", "")
	// addIndex("comment_images_versions", "image_ids", "") // postgres array

	addIndex("comment_parent_id_versions", "timestamp", "")
	addIndex("comment_parent_id_versions", "item_id", "")
	addIndex("comment_parent_id_versions", "value", "")

	addIndex("comment_creating_timestamp_versions", "timestamp", "")
	addIndex("comment_creating_timestamp_versions", "item_id", "")
	addIndex("comment_creating_timestamp_versions", "value", "")

	addIndex("comment_rating_versions", "timestamp", "")
	addIndex("comment_rating_versions", "item_id", "")
	addIndex("comment_rating_versions", "value", "")

	addIndex("comment_story_id_versions", "timestamp", "")
	addIndex("comment_story_id_versions", "item_id", "")
	addIndex("comment_story_id_versions", "value", "")

	addIndex("comment_user_id_versions", "timestamp", "")
	addIndex("comment_user_id_versions", "item_id", "")
	addIndex("comment_user_id_versions", "value", "")

	addIndex("comment_author_username_versions", "timestamp", "")
	addIndex("comment_author_username_versions", "item_id", "")
	addIndex("comment_author_username_versions", "value", "hash")

	addIndex("comment_is_hidden_versions", "timestamp", "")
	addIndex("comment_is_hidden_versions", "item_id", "")
	addIndex("comment_is_hidden_versions", "value", "")

	addIndex("comment_is_deleted_versions", "timestamp", "")
	addIndex("comment_is_deleted_versions", "item_id", "")
	addIndex("comment_is_deleted_versions", "value", "")

	addIndex("comment_is_author_community_moderator_versions", "timestamp", "")
	addIndex("comment_is_author_community_moderator_versions", "item_id", "")
	addIndex("comment_is_author_community_moderator_versions", "value", "")

	addIndex("comment_is_author_pikabu_team_versions", "timestamp", "")
	addIndex("comment_is_author_pikabu_team_versions", "item_id", "")
	addIndex("comment_is_author_pikabu_team_versions", "value", "")

	addIndex("comment_text_versions", "comment_id", "")
	addIndex("comment_text_versions", "timestamp", "")
	addIndex("comment_text_versions", "diffs", "hash")
}
