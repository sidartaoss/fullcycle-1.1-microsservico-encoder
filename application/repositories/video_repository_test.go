package repositories

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sidartaoss/fullcycle/encoder/domain"
	"github.com/sidartaoss/fullcycle/encoder/framework/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var resourceID = "ResourceID"
var filePath = "FilePath"

func TestVideoRepositoryDBInsert(t *testing.T) {
	// arrange
	db := database.NewDBTest()
	assert.NotNil(t, db)
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.ResourceID = resourceID
	video.FilePath = filePath
	video.CreatedAt = time.Now()

	repo := NewVideoRepositoryDB(db)
	assert.NotNil(t, repo)

	// act
	dbVideo, err := repo.Insert(video)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, dbVideo)
	assert.Equal(t, video.ID, dbVideo.ID)
	assert.Equal(t, video.ResourceID, dbVideo.ResourceID)
	assert.Equal(t, video.FilePath, dbVideo.FilePath)
	assert.Equal(t, video.CreatedAt, dbVideo.CreatedAt)

	v, err := repo.Find(video.ID)
	assert.Nil(t, err)
	require.NotEmpty(t, v.ID)
	require.Equal(t, video.ID, v.ID)
}
