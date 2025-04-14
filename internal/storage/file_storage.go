package storage

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"notes-api/internal/models"
	"os"
	"sync"
	"time"
)

type FileStorage struct {
	filePath string
	mu       sync.Mutex
}

func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{filePath: filePath}
}

func (fs *FileStorage) SaveAll(notes *[]models.Note) error {
	data, err := json.MarshalIndent(*notes, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.filePath, data, 0644)
}

func (fs *FileStorage) LoadAll() ([]models.Note, error) {
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Note{}, nil
		}
		return nil, err
	}

	var notes []models.Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

func (fs *FileStorage) GetAll() ([]models.Note, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	return fs.LoadAll()
}

func (fs *FileStorage) GetByID(id string) (*models.Note, error) {
	notes, err := fs.GetAll()
	if err != nil {
		return nil, err
	}

	for _, note := range notes {
		if note.ID == id {
			foundNote := note
			return &foundNote, nil
		}
	}

	return nil, errors.New("note not found")

}

func (fs *FileStorage) Create(note models.Note) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	notes, err := fs.LoadAll()
	if err != nil {
		return err
	}

	note.ID = GenerateID()
	note.CreatedAt = time.Now()

	notes = append(notes, note)
	return fs.SaveAll(&notes)
}

func (fs *FileStorage) Update(id string, updatedNote models.Note) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	notes, err := fs.LoadAll()
	if err != nil {
		return err
	}

	for i, note := range notes {
		if note.ID == id {
			updatedNote.ID = note.ID
			updatedNote.CreatedAt = note.CreatedAt

			notes[i] = updatedNote
			return fs.SaveAll(&notes)
		}
	}

	return errors.New("note not found")
}

func (fs *FileStorage) Delete(id string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	notes, err := fs.LoadAll()
	if err != nil {
		return err
	}

	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			return fs.SaveAll(&notes)
		}
	}

	return errors.New("note not found")
}

func GenerateID() string {
	return uuid.New().String()
}
