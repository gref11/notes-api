package storage

import (
	"encoding/json"
	"notes-api/internal/models"
	"os"
)

type FileStorage struct {
	filePath string
}

func NewFileStorage (filePath string) *FileStorage {
	return &FileStorage{filePath: filePath}
}

func (fs *FileStorage) SaveAll (notes *[]models.Note) error {
	data, err := json.MarshalIndent(notes, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.filePath, data, 0644)
}

func (fs *FileStorage) LoadAll () (notes []models.Note, err error) {
	
}