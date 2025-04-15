package storage

import (
	"errors"
	"notes-api/internal/models"
)

type Storage interface {
	GetAll() ([]models.Note, error)
	GetByID(id string) (*models.Note, error)
	Create(note models.Note) error
	Update(id string, note models.Note) error
	Delete(id string) error
}

var ErrNoteNotFound = errors.New("note not found")