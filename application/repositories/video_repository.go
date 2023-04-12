package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sidartaoss/fullcycle/encoder/domain"
)

type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDB struct {
	*gorm.DB
}

func NewVideoRepositoryDB(db *gorm.DB) *VideoRepositoryDB {
	return &VideoRepositoryDB{
		DB: db,
	}
}

func (r *VideoRepositoryDB) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.New().String()
	}
	err := r.DB.Create(video).Error
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (r *VideoRepositoryDB) Find(id string) (*domain.Video, error) {
	var video domain.Video
	err := r.DB.Preload("Jobs").First(&video, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	if video.ID == "" {
		return nil, errors.New("video does not exist")
	}
	return &video, nil
}
