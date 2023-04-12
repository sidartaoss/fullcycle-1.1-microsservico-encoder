package services

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {
	// arrange
	video, repository := prepare()

	videoService := NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repository

	err := videoService.Download("fullcyclesidartaoss")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = "fullcyclesidartaoss"
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + video.ID

	concurrency := 50
	doneUpload := make(chan string)

	// act
	go videoUpload.ProcessUpload(concurrency, doneUpload)
	result := <-doneUpload

	// assert
	require.NotNil(t, result)
	require.Equal(t, "upload completed", result)

	err = videoService.Finish()
	require.Nil(t, err)
}
