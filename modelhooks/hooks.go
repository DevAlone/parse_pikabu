package modelhooks

import (
	"github.com/DevAlone/parse_pikabu/models"
)

func HandleModelCreated(model interface{}, timestamp models.TimestampType) {
	// to be called from resultprocessor package
	modelCreatedEvents <- &ModelCreatedEvent{
		Data:      model,
		EventTime: timestamp,
	}
}

func HandleModelChanged(prev, curr interface{}, changeTime models.TimestampType) {
	// to be called from resultprocessor package
	modelChangedEvents <- &ModelChangedEvent{
		PrevData:  prev,
		CurrData:  curr,
		EventTime: changeTime,
	}
}

type ModelCreatedEvent struct {
	Data      interface{}
	EventTime models.TimestampType
}

type ModelChangedEvent struct {
	PrevData  interface{}
	CurrData  interface{}
	EventTime models.TimestampType
}

var modelCreatedEvents = make(chan *ModelCreatedEvent)
var modelChangedEvents = make(chan *ModelChangedEvent)

func RunModelHooksHandler() error {

	for {
		select {
		case modelCreatedEvent := <-modelCreatedEvents:
			err := handleModelCreatedEvent(modelCreatedEvent)
			if err != nil {
				return err
			}
		case modelChangedEvent := <-modelChangedEvents:
			err := handleModelChangedEvent(modelChangedEvent)
			if err != nil {
				return err
			}
		}
	}
}

func handleModelCreatedEvent(modelCreatedEvent *ModelCreatedEvent) error {
	// TODO: collect all errors and only then return
	switch data := modelCreatedEvent.Data.(type) {
	case models.PikabuUser:
		for _, handler := range pikabuUserCreatedEventHandlers {
			err := handler(&data, modelCreatedEvent.EventTime)
			if err != nil {
				return err
			}
		}
	case models.PikabuCommunity:
		for _, handler := range pikabuCommunityCreatedEventHandlers {
			err := handler(&data, modelCreatedEvent.EventTime)
			if err != nil {
				return err
			}
		}
	case models.PikabuStory:
		for _, handler := range pikabuStoryCreatedEventHandlers {
			err := handler(&data, modelCreatedEvent.EventTime)
			if err != nil {
				return err
			}
		}
	case models.PikabuComment:
		for _, handler := range pikabuCommentCreatedEventHandlers {
			err := handler(&data, modelCreatedEvent.EventTime)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func handleModelChangedEvent(modelChangedEvent *ModelChangedEvent) error {
	switch prevData := modelChangedEvent.PrevData.(type) {
	case models.PikabuUser:
		currData := modelChangedEvent.PrevData.(models.PikabuUser)
		for _, handler := range pikabuUserChangedEventHandlers {
			err := handler(&prevData, &currData, modelChangedEvent.EventTime)
			if err != nil {
				return err
			}
		}
	case models.PikabuCommunity:
		currData := modelChangedEvent.PrevData.(models.PikabuCommunity)
		for _, handler := range pikabuCommunityChangedEventHandlers {
			err := handler(&prevData, &currData, modelChangedEvent.EventTime)
			if err != nil {
				return err
			}
		}
	case models.PikabuStory:
		currData := modelChangedEvent.PrevData.(models.PikabuStory)
		for _, handler := range pikabuStoryChangedEventHandlers {
			err := handler(&prevData, &currData, modelChangedEvent.EventTime)
			if err != nil {
				return err
			}
		}
	case models.PikabuComment:
		currData := modelChangedEvent.PrevData.(models.PikabuComment)
		for _, handler := range pikabuCommentChangedEventHandlers {
			err := handler(&prevData, &currData, modelChangedEvent.EventTime)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
