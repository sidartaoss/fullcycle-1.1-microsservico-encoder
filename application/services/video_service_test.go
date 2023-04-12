package services

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sidartaoss/fullcycle/encoder/application/repositories"
	"github.com/sidartaoss/fullcycle/encoder/domain"
	"github.com/sidartaoss/fullcycle/encoder/framework/database"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepository) {
	db := database.NewDBTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "convite.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepositoryDB(db)
	_, err := repo.Insert(video)
	if err != nil {
		panic(err)
	}

	return video, repo
}

func TestVideoServiceDownload(t *testing.T) {
	// arrange
	video, repo := prepare()

	videoService := NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	// act & assert
	err := videoService.Download("fullcyclesidartaoss")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	err = videoService.Finish()
	require.Nil(t, err)
}
