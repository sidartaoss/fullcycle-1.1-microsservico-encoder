package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	// arrange
	video := NewVideo()

	// act
	err := video.Validate()

	// assert
	require.Error(t, err)
}

func TestVideoIdIsNotUUID(t *testing.T) {
	// arrange
	video := NewVideo()
	video.ID = "abc"
	video.ResourceID = "ResourceID"
	video.FilePath = "FilePath"
	video.CreatedAt = time.Now()

	// act
	err := video.Validate()

	// assert
	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	// arrange
	video := NewVideo()
	video.ID = uuid.New().String()
	video.ResourceID = "ResourceID"
	video.FilePath = "FilePath"
	video.CreatedAt = time.Now()

	// act
	err := video.Validate()

	// assert
	assert.Nil(t, err)
}
