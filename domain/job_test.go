package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	// arrange
	video := NewVideo()
	video.ID = uuid.New().String()
	video.ResourceID = "ResourceID"
	video.FilePath = "FilePath"
	video.CreatedAt = time.Now()

	outputBucketPath := "path"
	status := "Converted"

	// act
	job, err := NewJob(outputBucketPath, status, video)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, job)

}
