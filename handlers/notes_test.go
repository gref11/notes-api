package handlers

import (
	// "notes-api/internal/storage"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"notes-api/internal/models"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) GetAll() ([]models.Note, error) {
	args := m.Called()
	return args.Get(0).([]models.Note), args.Error(1)
}

func (m *StorageMock) GetByID(id string) (*models.Note, error) {
	args := m.Called()
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *StorageMock) Create(note models.Note) error {
	args := m.Called()
	return args.Error(0)
}

func (m *StorageMock) Update(id string, note models.Note) error {
	args := m.Called()
	return args.Error(0)
}

func (m *StorageMock) Delete(id string) error {
	args := m.Called()
	return args.Error(0)
}

func TestGetAllNotes(t *testing.T) {
	mockStorage := new(StorageMock)
	expectedNotes := []models.Note{
		{ID: "00001", Title: "Title", Content: "Content", CreatedAt: time.Date(2025, time.April, 1, 12, 0, 0, 0, time.UTC)},
	}
	mockStorage.On("GetAll").Return(expectedNotes, nil)

	h := NewNotesHandler(mockStorage)

	req, err := http.NewRequest("GET", "/notes", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/notes", h.GetAllNotes).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.Note
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedNotes, response)
	mockStorage.AssertExpectations(t)
}

func TestGetNoteByID(t *testing.T) {
	mockStorage := new(StorageMock)
	expectedNote := &models.Note{
		ID:        "00001",
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Date(2025, time.April, 1, 12, 0, 0, 0, time.UTC),
	}
	mockStorage.On("GetByID", mock.Anything).Return(expectedNote, nil)

	h := NewNotesHandler(mockStorage)

	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/notes/00001", nil)
	assert.NoError(t, err)

	router := mux.NewRouter()
	router.HandleFunc("/notes/{id}", h.GetNoteByID).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Note
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, *expectedNote, response)
	mockStorage.AssertExpectations(t)
}

func TestCreateNote(t *testing.T) {
	MockStorage := new(StorageMock)
	newNote := models.Note{
		ID:        "00002",
		Title:     "Note title",
		Content:   "Note content",
		CreatedAt: time.Date(2025, time.April, 23, 12, 0, 0, 0, time.UTC),
	}
	MockStorage.On("Create", mock.Anything).Return(nil)

	h := NewNotesHandler(MockStorage)

	body, err := json.Marshal(newNote)
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", "/notes/00002", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/notes/{id}", h.CreateNote).Methods("POST")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	MockStorage.AssertExpectations(t)
}
