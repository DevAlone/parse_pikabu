package old_models

type Image struct {
	Id                  uint64
	SmallURL            string `sql:",notnull,unique:comment_image__small_url__large_url__animation_base_url__animation_preview_url"`
	LargeURL            string `sql:",notnull,unique:comment_image__small_url__large_url__animation_base_url__animation_preview_url"`
	AnimationBaseURL    string `sql:",notnull,unique:comment_image__small_url__large_url__animation_base_url__animation_preview_url"`
	AnimationPreviewURL string `sql:",notnull,unique:comment_image__small_url__large_url__animation_base_url__animation_preview_url"`
	AnimationFormats    map[string]int
	// width and height can be null
	Width  int32
	Height int32
}

func init() {
	/*
		addIndex("images", "small_url", "hash")
		addIndex("images", "large_url", "hash")
	*/
	/*
		addIndex("images", "width", "")
		addIndex("images", "height", "")
	*/
}
