package modelhooks

// TODO: make this file generated

import (
	"github.com/DevAlone/parse_pikabu/models"
)

type PikabuUserCreatedEventHandler func(
	data *models.PikabuUser,
	eventTime models.TimestampType,
) error
type PikabuCommunityCreatedEventHandler func(
	data *models.PikabuCommunity,
	eventTime models.TimestampType,
) error
type PikabuStoryCreatedEventHandler func(
	data *models.PikabuStory,
	eventTime models.TimestampType,
) error
type PikabuCommentCreatedEventHandler func(
	data *models.PikabuComment,
	eventTime models.TimestampType,
) error

var pikabuUserCreatedEventHandlers = []PikabuUserCreatedEventHandler{}
var pikabuCommunityCreatedEventHandlers = []PikabuCommunityCreatedEventHandler{}
var pikabuStoryCreatedEventHandlers = []PikabuStoryCreatedEventHandler{}
var pikabuCommentCreatedEventHandlers = []PikabuCommentCreatedEventHandler{}

func SubscribeToPikabuUserCreatedEvent(handler PikabuUserCreatedEventHandler) {
	pikabuUserCreatedEventHandlers = append(pikabuUserCreatedEventHandlers, handler)
}
func SubscribeToPikabuCommunityCreatedEvent(handler PikabuCommunityCreatedEventHandler) {
	pikabuCommunityCreatedEventHandlers = append(pikabuCommunityCreatedEventHandlers, handler)
}
func SubscribeToPikabuStoryCreatedEvent(handler PikabuStoryCreatedEventHandler) {
	pikabuStoryCreatedEventHandlers = append(pikabuStoryCreatedEventHandlers, handler)
}
func SubscribeToPikabuCommentCreatedEvent(handler PikabuCommentCreatedEventHandler) {
	pikabuCommentCreatedEventHandlers = append(pikabuCommentCreatedEventHandlers, handler)
}

type PikabuUserChangedEventHandler func(
	prevData *models.PikabuUser,
	currData *models.PikabuUser,
	eventTime models.TimestampType,
) error
type PikabuCommunityChangedEventHandler func(
	prevData *models.PikabuCommunity,
	currData *models.PikabuCommunity,
	eventTime models.TimestampType,
) error
type PikabuStoryChangedEventHandler func(
	prevData *models.PikabuStory,
	currData *models.PikabuStory,
	eventTime models.TimestampType,
) error
type PikabuCommentChangedEventHandler func(
	prevData *models.PikabuComment,
	currData *models.PikabuComment,
	eventTime models.TimestampType,
) error

var pikabuUserChangedEventHandlers = []PikabuUserChangedEventHandler{}
var pikabuCommunityChangedEventHandlers = []PikabuCommunityChangedEventHandler{}
var pikabuStoryChangedEventHandlers = []PikabuStoryChangedEventHandler{}
var pikabuCommentChangedEventHandlers = []PikabuCommentChangedEventHandler{}

func SubscribeToPikabuUserChangedEvent(handler PikabuUserChangedEventHandler) {
	pikabuUserChangedEventHandlers = append(pikabuUserChangedEventHandlers, handler)
}
func SubscribeToPikabuCommunityChangedEvent(handler PikabuCommunityChangedEventHandler) {
	pikabuCommunityChangedEventHandlers = append(pikabuCommunityChangedEventHandlers, handler)
}
func SubscribeToPikabuStoryChangedEvent(handler PikabuStoryChangedEventHandler) {
	pikabuStoryChangedEventHandlers = append(pikabuStoryChangedEventHandlers, handler)
}
func SubscribeToPikabuCommentChangedEvent(handler PikabuCommentChangedEventHandler) {
	pikabuCommentChangedEventHandlers = append(pikabuCommentChangedEventHandlers, handler)
}
