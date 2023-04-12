package repositories

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/sidartaoss/fullcycle/encoder/domain"
)

type JobRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDB struct {
	*gorm.DB
}

func NewJobRepositoryDB(db *gorm.DB) *JobRepositoryDB {
	return &JobRepositoryDB{
		DB: db,
	}
}

func (r *JobRepositoryDB) Insert(job *domain.Job) (*domain.Job, error) {
	err := r.DB.Create(job).Error
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r *JobRepositoryDB) Find(id string) (*domain.Job, error) {
	var job domain.Job
	err := r.DB.Preload("Video").First(&job, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	if job.ID == "" {
		return nil, errors.New("job does not exist")
	}
	return &job, nil
}

func (r *JobRepositoryDB) Update(job *domain.Job) (*domain.Job, error) {
	err := r.DB.Save(job).Error
	if err != nil {
		return nil, err
	}
	return job, nil
}
