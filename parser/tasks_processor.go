package parser

import (
	"time"

	"github.com/DevAlone/parse_pikabu/core/taskmanager"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"
)

// Loop - parser's loop
func (p *Parser) Loop() {
	for {
		err := p.ListenForTasks()
		if err != nil {
			p.handleError(err)
			time.Sleep(time.Duration(p.Config.WaitAfterErrorSeconds) * time.Second)
		}
	}
}

// ListenForTasks -
func (p *Parser) ListenForTasks() error {
	defer func() {
		if r := recover(); r != nil {
			p.handleError(errors.Errorf("panic: %v", r))
		}
	}()

	for {
		task, err := taskmanager.GetParserTask()
		if err != nil {
			return err
		}
		switch t := task.(type) {
		case *models.ParseUserTask:
			return p.processParseUserTask(t)
		case *models.ParseStoryTask:
			return p.processParseStoryTask(t)
		case *models.SimpleTask:
			switch t.Name {
			case "parse_communities_pages":
				return p.processParseCommunitiesPagesTask()
			default:
				return errors.Errorf("unknown type of simple task %v", t)
			}
		default:
			return errors.Errorf("Parser got unknown task %v", t)
		}
	}
}
