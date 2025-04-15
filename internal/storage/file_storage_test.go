package storage

import (
	"notes-api/internal/models"
	"os"
	"testing"
)

func TestFileStorageCRUD(t *testing.T) {
	path := "storage_test.json"
	storage := NewFileStorage(path)
	defer os.Remove(path)

	testNote := models.Note{
		Title:   "Title",
		Content: "Content",
	}

	if err := storage.Create(testNote); err != nil {
		t.Errorf("Failed note create: %v", err)
	}

	notes, err := storage.GetAll()
	if err != nil {
		t.Errorf("Failed notes get: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("Failed notes get: expected 1 note, got %d", len(notes))
	}

	noteID := notes[0].ID
	note, err := storage.GetByID(noteID)
	if err != nil {
		t.Errorf("Failed note get by ID: %v", err)
	}
	if note.Title != testNote.Title {
		t.Errorf("Title expected: %q, got: %q", testNote.Title, note.Title)
	}

	updatedNote := models.Note{
		Title:   "Updated title",
		Content: "Updated content",
	}
	if err := storage.Update(noteID, updatedNote); err != nil {
		t.Errorf("Failed note update: %v", err)
	}
	note, err = storage.GetByID(noteID)
	if err != nil {
		t.Errorf("Failed note update: %v", err)
	}
	if note.Title != updatedNote.Title {
		t.Errorf("Failed note update: title excpected: %q, got: %q", updatedNote.Title, note.Title)
	}

	if err := storage.Delete(noteID); err != nil {
		t.Errorf("Failed note delete: %v", err)
	}
	if _, err := storage.GetByID(noteID); err == nil {
		t.Errorf("Failed note delete: note was not deleted")
	}
}
