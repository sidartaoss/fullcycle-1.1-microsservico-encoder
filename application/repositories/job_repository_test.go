package repositories

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sidartaoss/fullcycle/encoder/domain"
	"github.com/sidartaoss/fullcycle/encoder/framework/database"
	"github.com/stretchr/testify/assert"
)

func TestJobRepositoryDBInsert(t *testing.T) {
	// arrange
	db := database.NewDBTest()
	assert.NotNil(t, db)
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.ResourceID = "ResourceID"
	video.FilePath = "FilePath"
	video.CreatedAt = time.Now()

	videoRepository := NewVideoRepositoryDB(db)
	assert.NotNil(t, videoRepository)

	dbVideo, err := videoRepository.Insert(video)
	assert.Nil(t, err)
	assert.NotNil(t, dbVideo)

	job, err := domain.NewJob("path", "Converted", video)
	assert.Nil(t, err)
	assert.NotNil(t, job)

	jobRepository := NewJobRepositoryDB(db)
	assert.NotNil(t, jobRepository)

	// act
	dbJob, err := jobRepository.Insert(job)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, dbJob)

	j, err := jobRepository.Find(dbJob.ID)
	assert.Nil(t, err)
	assert.NotNil(t, j)
	assert.Equal(t, job.ID, j.ID)
	assert.NotNil(t, j.Video)
	assert.Equal(t, video.ID, j.VideoID)
}

func TestJobRepositoryDBUpdate(t *testing.T) {
	// arrange
	db := database.NewDBTest()
	assert.NotNil(t, db)
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.ResourceID = "ResourceID"
	video.FilePath = "FilePath"
	video.CreatedAt = time.Now()

	videoRepository := NewVideoRepositoryDB(db)
	assert.NotNil(t, videoRepository)

	dbVideo, err := videoRepository.Insert(video)
	assert.Nil(t, err)
	assert.NotNil(t, dbVideo)

	job, err := domain.NewJob("path", "Converted", video)
	assert.Nil(t, err)
	assert.NotNil(t, job)

	jobRepository := NewJobRepositoryDB(db)
	assert.NotNil(t, jobRepository)

	dbJob, err := jobRepository.Insert(job)
	assert.Nil(t, err)
	assert.NotNil(t, dbJob)

	job.Status = "Completed"

	// act
	dbJob, err = jobRepository.Update(job)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, dbJob)

	j, err := jobRepository.Find(dbJob.ID)
	assert.Nil(t, err)
	assert.NotNil(t, j)
	assert.Equal(t, job.ID, j.ID)
	assert.NotNil(t, j.Video)
	assert.Equal(t, video.ID, j.VideoID)
	assert.Equal(t, job.Status, j.Status)
}
