package tests

import (
	"testing"

	"bitbucket.org/d3dev/parse_pikabu/core/config"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/d3dev/parse_pikabu/models"

	"bitbucket.org/d3dev/parse_pikabu/parser"
)

func TestParseUserByIdTaskProcessing(t *testing.T) {
	initLogs()
	config.Settings.Debug = true

	conf, err := parser.NewParserConfigFromBytes([]byte(
		`{"ProxyProviderAPIURL": "https://eivailohciihi4uquapach7abei9iesh.d3d.info"}`,
	))
	if err != nil {
		t.Fatal(err)
	}

	p, err := parser.NewParser(conf)
	if err != nil {
		t.Fatal(err)
	}

	/*
		res, err := p.ProcessParseUserTaskById(models.ParseUserTask{
			PikabuId:       1,
			IsDone:         false,
			IsTaken:        true,
			AddedTimestamp: 0,
			Username:       "",
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "admin", res.User.Username)
	*/

	res, err := p.ProcessParseUserTaskById(models.ParseUserTask{
		// PikabuId:       2159282,
		PikabuId:       4,
		IsDone:         false,
		IsTaken:        true,
		AddedTimestamp: 0,
		Username:       "",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "admin", res.User.Username)
}
