package models

type PikabuCommunity struct {
	PikabuId            uint64   `sql:",pk"`
	Name                string   `sql:",notnull" gen_versions:""`
	LinkName            string   `sql:",notnull" gen_versions:""`
	URL                 string   `sql:",notnull" gen_versions:""`
	AvatarURL           string   `sql:",notnull" gen_versions:""`
	BackgroundImageURL  string   `sql:",notnull" gen_versions:""`
	Tags                []string `sql:",notnull,array" gen_versions:""`
	NumberOfStories     int32    `sql:",notnull" gen_versions:""`
	NumberOfSubscribers int32    `sql:",notnull" gen_versions:""`
	Description         string   `sql:",notnull" gen_versions:""`
	Rules               string   `sql:",notnull" gen_versions:""`
	Restrictions        string   `sql:",notnull" gen_versions:""`
	AdminId             uint64   `sql:",notnull" gen_versions:""`
	ModeratorIds        []uint64 `sql:",notnull" gen_versions:""`

	AddedTimestamp      TimestampType `sql:",notnull"`
	LastUpdateTimestamp TimestampType `sql:",notnull"`
}

func init() {
	for _, item := range []interface{}{
		&PikabuCommunity{},
	} {
		Tables = append(Tables, item)
	}

	/*
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
	*/
}
