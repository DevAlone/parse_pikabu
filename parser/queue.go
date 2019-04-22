package parser

import (
	"reflect"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/globals"

	"bitbucket.org/d3dev/parse_pikabu/models"
)

// PutResultsToQueue - handles parser's results
func (p *Parser) PutResultsToQueue(routingKey string, result interface{}) error {
	numberOfResults := 0
	resultType := reflect.TypeOf(result)
	switch resultType.Kind() {
	case reflect.Slice, reflect.Array:
		numberOfResults = reflect.ValueOf(result).Len()
	default:
		result = []interface{}{result}
		numberOfResults = 1
	}

	var pr models.ParserResult
	pr.ParsingTimestamp = models.TimestampType(time.Now().Unix())
	pr.ParserId = "d3dev/" + p.Config.ParserId
	pr.NumberOfResults = numberOfResults
	pr.Results = result
	globals.ParserResults <- &pr

	return nil
}