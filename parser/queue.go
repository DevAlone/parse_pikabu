package parser

import (
	"reflect"
	"time"

	"github.com/DevAlone/parse_pikabu/globals"

	"github.com/DevAlone/parse_pikabu/models"
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
	pr.ParserID = "d3dev/" + p.Config.ParserID
	pr.NumberOfResults = numberOfResults
	pr.Results = result
	globals.ParserResults <- &pr

	return nil
}
