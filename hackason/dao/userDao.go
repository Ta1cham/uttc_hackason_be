package dao

import (
	"database/sql"
	"errors"
	"log"
	"uttc_hackason_be/model"
)

type UserDao struct {
	DB *sql.DB
}

func NewUserDao(db *sql.DB) *UserDao {
	dao := &UserDao{DB: db}
	return dao
}

func (dao *UserDao) RegisterUser(user *model.UserInfoForHTTPPOST) error {
	db := dao.DB
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO user (id, name) VALUES (?, ?)", user.ID, user.Name); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (dao *UserDao) SearchUser(id string) (model.UserResForHTTPGET, error) {
	db := dao.DB
	var u model.UserResForHTTPGET
	tx, err := db.Begin()
	if err != nil {
		log.Printf("failied to begin transaction")
		return u, err
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

	query := "SELECT id, name FROM user WHERE id = ?"
	err = tx.QueryRow(query, id).Scan(&u.ID, &u.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, nil
		}
		return u, err
	}

	return u, nil
}
