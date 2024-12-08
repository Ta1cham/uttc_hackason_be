package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"uttc_hackason_be/model"
	"uttc_hackason_be/usecase"
)

type NoteController struct {
	NoteUseCase *usecase.NoteUseCase
}

func NewNoteController(db *sql.DB) *NoteController {
	noteUseCase := usecase.NewNoteUseCase(db)
	return &NoteController{NoteUseCase: noteUseCase}
}

func (nc *NoteController) addNote(w http.ResponseWriter, r *http.Request) {
	var content *model.ContentData

	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(content)
	note, err := nc.NoteUseCase.AddNote(content.Pid, content.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	// noteの送信
	response := map[string]string{
		"note": note,
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (uc *NoteController) RegiterRoutes(r *mux.Router) {
	r.HandleFunc("/note", uc.addNote).Methods("POST")
}
