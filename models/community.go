package models

type PikabuCommunity struct {
	PikabuId            uint64   `sql:",pk" json:"pikabu_id" api:"ordering,filter"`
	Name                string   `sql:",notnull" gen_versions:"" json:"name" api:"ordering,filter"`
	LinkName            string   `sql:",notnull" gen_versions:"" json:"link_name" api:"filter"`
	URL                 string   `sql:",notnull" gen_versions:"" json:"url" api:"filter"`
	AvatarURL           string   `sql:",notnull" gen_versions:"" json:"avatar_url" api:"filter"`
	BackgroundImageURL  string   `sql:",notnull" gen_versions:"" json:"background_image_url" api:"filter"`
	Tags                []string `sql:",notnull,array" gen_versions:"" json:"tags" api:"filter"`
	NumberOfStories     int32    `sql:",notnull" gen_versions:"" json:"number_of_stories" api:"ordering,filter"`
	NumberOfSubscribers int32    `sql:",notnull" gen_versions:"" json:"number_of_subscribers" api:"ordering,filter"`
	Description         string   `sql:",notnull" gen_versions:"" json:"description" api:"filter"`
	Rules               string   `sql:",notnull" gen_versions:"" json:"rules" api:"filter"`
	Restrictions        string   `sql:",notnull" gen_versions:"" json:"restrictions" api:"filter"`
	AdminId             uint64   `sql:",notnull" gen_versions:"" json:"admin_id" api:"filter"`
	ModeratorIds        []uint64 `sql:",notnull" gen_versions:"" json:"moderator_ids" api:"filter"`

	AddedTimestamp      TimestampType `sql:",notnull" json:"added_timestamp" api:"ordering,filter"`
	LastUpdateTimestamp TimestampType `sql:",notnull" json:"last_update_timestamp" api:"ordering,filter"`
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
