package models

import (
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

type ParserBaseResult struct  {
	ParsingTimestamp TimestampType `json:"parsing_timestamp"`
	ParserId         string        `json:"parser_id"`
	NumberOfResults  int           `json:"number_of_results"`
}

type ParserResult struct {
	ParserBaseResult
	Results interface{} `json:"results"`
}

type ParserUserProfileResult struct {
	ParserBaseResult
	Results []struct {
		User *pikago_models.UserProfile `json:"user"`
	} `json:"results"`
}

type ParserCommunitiesPageResult struct {
	ParserBaseResult
	Results []pikago_models.CommunitiesPage `json:"results"`
}

type ParserUserProfileNotFoundResultData struct {
	PikabuId uint64 `json:"pikabu_id"`
	Username string `json:"username"`
	PikabuError interface{} `json:"pikabu_error"`
}

type ParserUserProfileNotFoundResult struct {
	ParserBaseResult
	Results []ParserUserProfileNotFoundResultData `json:"results"`
}
