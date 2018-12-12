package models

type Community struct {
	TableName   struct{} `sql:"communities_app_community"`
	Id          uint64
	UrlName     string
	Name        string
	Description string
	// TODO: add history
	AvatarURL           string
	BackgroundImageURL  string
	SubscribersCount    int32
	StoriesCount        int32
	LastUpdateTimestamp int64
}

type CommunityCountersEntry struct {
	TableName struct{} `sql:"communities_app_communitycountersentry"`

	Id               uint64
	Timestamp        int64
	CommunityId      uint64
	SubscribersCount int32
	StoriesCount     int32
}

func init() {
	addIndex("communities_app_community", "url_name", "hash")
	addUniqueIndex("communities_app_community", "url_name", "")
	addIndex("communities_app_community", "name", "")
	// addIndex("communities_app_community", "name", "hash")
	// addIndex("communities_app_community", "description", "hash")
	// addIndex("communities_app_community", "avatar_url", "hash")
	// addIndex("communities_app_community", "background_image_url", "hash")
	addIndex("communities_app_community", "subscribers_count", "")
	addIndex("communities_app_community", "stories_count", "")
	addIndex("communities_app_community", "last_update_timestamp", "")

	addIndex("communities_app_communitycountersentry", "timestamp", "")
	addIndex("communities_app_communitycountersentry", "community_id", "")
	// addIndex("communities_app_communitycountersentry", "subscribers_count", "")
	// addIndex("communities_app_communitycountersentry", "stories_count", "")

	addUniqueIndex("communities_app_communitycountersentry", []string{"timestamp", "community_id"}, "")
}
