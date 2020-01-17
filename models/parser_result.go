package models

import (
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

// ParserBaseResult -
type ParserBaseResult struct {
	ParsingTimestamp TimestampType `json:"parsing_timestamp"`
	ParserID         string        `json:"parser_id"`
	NumberOfResults  int           `json:"number_of_results"`
	// TypeOfResult     string        `json:"type_of_result"`
}

type ParserResult struct {
	ParserBaseResult
	Results interface{} `json:"results"`
}

type ParserUserProfileResultData struct {
	User *pikago_models.UserProfile `json:"user"`
}

type ParserUserProfileResult struct {
	ParserBaseResult
	Results []ParserUserProfileResultData `json:"results"`
}

type ParserCommunitiesPageResult struct {
	ParserBaseResult
	Results []pikago_models.CommunitiesPage `json:"results"`
}

type ParserUserProfileNotFoundResultData struct {
	PikabuID    uint64      `json:"pikabu_id"`
	Username    string      `json:"username"`
	PikabuError interface{} `json:"pikabu_error"`
}

type ParserUserProfileNotFoundResult struct {
	ParserBaseResult
	Results []ParserUserProfileNotFoundResultData `json:"results"`
}

type ParserStoryNotFoundResultData struct {
	PikabuID    uint64      `json:"pikabu_id"`
	PikabuError interface{} `json:"pikabu_error"`
}

type ParserStoryHiddenInAPIResultData struct {
	PikabuID    uint64      `json:"pikabu_id"`
	PikabuError interface{} `json:"pikabu_error"`
}
