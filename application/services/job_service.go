package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/sidartaoss/fullcycle/encoder/application/repositories"
	"github.com/sidartaoss/fullcycle/encoder/domain"
)

type JobService struct {
	*domain.Job
	repositories.JobRepository
	VideoService
}

func NewJobService() *JobService {
	return &JobService{}
}

func (s *JobService) Start() error {

	err := s.ChangeJobStatus("DOWNLOADING")
	if err != nil {
		return s.FailJob(err)
	}

	err = s.VideoService.Download(os.Getenv("inputBucketName"))
	if err != nil {
		return s.FailJob(err)
	}

	err = s.ChangeJobStatus("FRAGMENTING")
	if err != nil {
		return s.FailJob(err)
	}

	err = s.VideoService.Fragment()
	if err != nil {
		return s.FailJob(err)
	}

	err = s.ChangeJobStatus("ENCODING")
	if err != nil {
		return s.FailJob(err)
	}

	err = s.VideoService.Encode()
	if err != nil {
		return s.FailJob(err)
	}

	err = s.PerformUpload()
	if err != nil {
		return s.FailJob(err)
	}

	err = s.ChangeJobStatus("FINISHING")
	if err != nil {
		return s.FailJob(err)
	}

	err = s.VideoService.Finish()
	if err != nil {
		return s.FailJob(err)
	}

	err = s.ChangeJobStatus("COMPLETED")
	if err != nil {
		return s.FailJob(err)
	}

	return nil
}

func (s *JobService) PerformUpload() error {

	err := s.ChangeJobStatus("UPLOADING")
	if err != nil {
		return s.FailJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + s.VideoService.Video.ID
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	if err != nil {
		return s.FailJob(err)
	}

	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	uploadResult := <-doneUpload

	if uploadResult != "upload completed" {
		return s.FailJob(errors.New(uploadResult))
	}

	return nil
}

func (s *JobService) ChangeJobStatus(status string) error {
	s.Job.Status = status
	_, err := s.JobRepository.Update(s.Job)
	if err != nil {
		return s.FailJob(err)
	}
	return nil
}

func (j *JobService) FailJob(fail error) error {
	j.Job.Status = "FAILED"
	j.Job.Error = fail.Error()

	_, err := j.JobRepository.Update(j.Job)

	if err != nil {
		return err
	}

	return fail
}
