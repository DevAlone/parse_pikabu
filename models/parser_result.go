package models

import (
	"gogsweb.2-47.ru/d3dev/pikago"
)

type ParserBaseResult struct {
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
		User *pikago.UserProfile `json:"user"`
	} `json:"results"`
}

type ParserCommunitiesPageResult struct {
	ParserBaseResult
	Results []pikago.CommunitiesPage `json:"results"`
}
