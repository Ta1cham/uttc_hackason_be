package dao

import (
	"database/sql"
	"log"
)

type NoteDao struct {
	DB *sql.DB
}

func NewNoteDao(db *sql.DB) *NoteDao {
	return &NoteDao{DB: db}
}

func (dao *NoteDao) AddNote(pid string, note string) error {
	db := dao.DB
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Fatalf("Transaction failed: %v", p)
		} else if err != nil {
			tx.Rollback()
			log.Fatalf("Transaction rolled back due to error: %v", err)
		} else {
			commitErr := tx.Commit()
			if err != nil {
				err = commitErr
				log.Fatalf("Failed to commit transaction: %v", commitErr)
			}
		}
	}()

	query := "INSERT INTO notes (tweet_id, note) VALUES (?, ?)"
	_, err = tx.Exec(query, pid, note)
	if err != nil {
		return err
	}
	return nil
}
