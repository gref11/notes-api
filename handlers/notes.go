package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"notes-api/internal/models"
	"notes-api/internal/storage"
)

type NotesHandler struct {
	storage storage.Storage
}

func NewNotesHandler(storage storage.Storage) *NotesHandler {
	return &NotesHandler{storage: storage}
}

func (h *NotesHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	notes, err := h.storage.GetAll()
	if err != nil {
		log.Printf("Storage error: %v", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *NotesHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, `{"error": "Bad request: id cannot be empty"}`, http.StatusBadRequest)
		return
	}


	note, err := h.storage.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoteNotFound):
			http.Error(w, `{"error": "Note not found"}`, http.StatusNotFound)
		default:
			log.Printf("Storage error: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *NotesHandler) CreateNote (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, `{"error": "Bad request: id cannot be empty"}`, http.StatusBadRequest)
		return
	}

	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, `{"error": "Bad request"}`, http.StatusBadRequest)
		return
	}

	if note.Title == "" {
		http.Error(w, `{"error": "Bad request: title cannot be empty"}`, http.StatusBadRequest)
		return
	}
	if note.Content == "" {
		http.Error(w, `{"error": "Bad request: content cannot be empty"}`, http.StatusBadRequest)
		return
	}

	if err := h.storage.Create(note); err != nil {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NotesHandler) UpdateNote (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, `{"error": "Bad request: id cannot be empty"}`, http.StatusBadRequest)
		return
	}

	var note models.Note

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, `{"error": "Bad request"}`, http.StatusBadRequest)
		return
	}

	if note.Title == "" {
		http.Error(w, `{"error": "Bad request: title cannot be empty"}`, http.StatusBadRequest)
		return
	}
	if note.Content == "" {
		http.Error(w, `{"error": "Bad request: content cannot be empty"}`, http.StatusBadRequest)
		return
	}

	if err := h.storage.Update(id, note); err != nil {
		switch {
		case errors.Is(err, storage.ErrNoteNotFound):
			http.Error(w, `{"error": "Note not found"}`, http.StatusNotFound)
		default:
			log.Printf("Storage error: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NotesHandler) DeleteNote (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, `{"error": "Bad request: id cannot be empty"}`, http.StatusBadRequest)
		return
	}

	if err := h.storage.Delete(id); err != nil {
		switch {
		case errors.Is(err, storage.ErrNoteNotFound):
			http.Error(w, `{"error": "Note not found"}`, http.StatusNotFound)
		default:
			log.Printf("Storage error: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}