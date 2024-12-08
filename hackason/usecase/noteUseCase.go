package usecase

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"uttc_hackason_be/dao"
	"uttc_hackason_be/model"
)

type NoteUseCase struct {
	NoteDao *dao.NoteDao
}

func NewNoteUseCase(db *sql.DB) *NoteUseCase {
	noteDao := dao.NewNoteDao(db)
	return &NoteUseCase{NoteDao: noteDao}
}

func genNote(content string) (string, error) {
	const location = "asia-northeast1"
	const modelName = "gemini-1.5-flash-002"
	//projectId := os.Getenv("GCP_PROJECT_ID")
	projectId := "term6-taichi-ogawa"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectId, location)
	if err != nil {
		return "", err
	}

	gemini := client.GenerativeModel(modelName)
	promptText := "以下の情報の信憑性は? " + content
	prompt := genai.Text(promptText)
	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		return "", err
	}
	if resp == nil {
		log.Printf("no response")
		return "", errors.New("no response")
	}
	rb, err := json.MarshalIndent(resp, "", "  ")
	jsonData := string(rb)

	var response model.Response
	err = json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", nil
	}
	note := response.Candidates[0].Content.Parts[0]
	cleanedNote := strings.ReplaceAll(note, "**", "")
	return cleanedNote, nil
}

func (nc *NoteUseCase) AddNote(pid string, content string) (string, error) {
	note, err := genNote(content)
	if err != nil {
		return "", err
	}
	if note == "" {
		return "", errors.New("empty note")
	}
	err = nc.NoteDao.AddNote(pid, note)
	if err != nil {
		return "", err
	}
	return note, nil
}
