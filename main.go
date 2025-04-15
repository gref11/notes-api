package main

import (
	"log"
	"net/http"
	"notes-api/handlers"
	"notes-api/internal/storage"
	"path/filepath"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	path := filepath.Join(".", "data", "notes.json")

	storage := storage.NewFileStorage(path)
	h := handlers.NewNotesHandler(storage)

	r.HandleFunc("/notes", h.GetAllNotes).Methods("GET")
	r.HandleFunc("/notes/{id}", h.GetNoteByID).Methods("GET")
	r.HandleFunc("/notes", h.CreateNote).Methods("POST")
	r.HandleFunc("/notes/{id}", h.UpdateNote).Methods("PUT")
	r.HandleFunc("/notes/{id}", h.DeleteNote).Methods("DELETE")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
