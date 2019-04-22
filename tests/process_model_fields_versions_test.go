package tests

import (
	"testing"

	"bitbucket.org/d3dev/parse_pikabu/core/resultsprocessor"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/stretchr/testify/assert"
)

func TestProcessModelFieldsVersions(t *testing.T) {
	oldModelPtr := &models.PikabuStory{
		ContentBlocks: []models.PikabuStoryBlock{},
	}
	newModelPtr := &models.PikabuStory{
		ContentBlocks: []models.PikabuStoryBlock{},
	}

	changedFields, err := resultsprocessor.GetModelsChangedFields(oldModelPtr, newModelPtr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, len(changedFields))

	///

	oldModelPtr = &models.PikabuStory{
		ContentBlocks: []models.PikabuStoryBlock{
			models.PikabuStoryBlock{
				Type: "type",
				Data: "some kind of data",
			},
			models.PikabuStoryBlock{
				Type: "type2",
				Data: "some other kind of data",
			},
		},
	}
	newModelPtr = &models.PikabuStory{
		ContentBlocks: []models.PikabuStoryBlock{
			models.PikabuStoryBlock{
				Type: "type",
				Data: "some kind of data",
			},
			models.PikabuStoryBlock{
				Type: "type2",
				Data: "some other kind of data",
			},
		},
	}

	changedFields, err = resultsprocessor.GetModelsChangedFields(oldModelPtr, newModelPtr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, len(changedFields))

	///

	oldModelPtr = &models.PikabuStory{
		ContentBlocks: []models.PikabuStoryBlock{
			models.PikabuStoryBlock{
				Type: "type",
				Data: "some kind of data",
			},
			models.PikabuStoryBlock{
				Type: "type2",
				Data: "some other kind of data",
			},
		},
	}
	newModelPtr = &models.PikabuStory{
		ContentBlocks: []models.PikabuStoryBlock{
			models.PikabuStoryBlock{
				Type: "asdfasdfasdfasdf",
				Data: "some kind of data",
			},
			models.PikabuStoryBlock{
				Type: "type2",
				Data: "some other kind of data",
			},
		},
	}

	changedFields, err = resultsprocessor.GetModelsChangedFields(oldModelPtr, newModelPtr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(changedFields))
}
