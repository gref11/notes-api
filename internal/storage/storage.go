package storage

import (
	"notes-api/internal/models"
)

type Storage interface {
	GetAll() ([]models.Note, error)
	GetByID(id string) (*models.Note, error)
	Create(note models.Note) error
	Update(id string, note models.Note) error
	Delete(id string) error
}
