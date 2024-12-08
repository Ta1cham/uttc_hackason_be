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
		log.Printf("Failed to begin transaction")
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
			if commitErr != nil {
				err = commitErr
				log.Fatalf("Failed to commit transaction: %v", commitErr)
			}
		}
	}()

	query := "SELECT id, name, bio, image FROM user WHERE id = ?"
	var bio sql.NullString
	var image sql.NullString

	err = tx.QueryRow(query, id).Scan(&u.ID, &u.Name, &bio, &image)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, nil
		}
		return u, err
	}

	// NULLチェックして構造体に代入
	u.Bio = bio.String
	if !bio.Valid {
		u.Bio = "" // NULLなら空文字列にする
	}

	u.Image = image.String
	if !image.Valid {
		u.Image = "" // NULLなら空文字列にする
	}

	return u, nil
}

func (dao *UserDao) EditProfile(info *model.EditInfoForHTTPPOST) error {
	db := dao.DB
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	query := `
		UPDATE user
		SET name = ?, bio = ?, image = ?
		WHERE id = ?
		`

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Fatalf("Transaction failed: %v", p)
		} else if err != nil {
			tx.Rollback()
			log.Fatalf("Transaction rolled back due to error: %v", err)
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				err = commitErr
				log.Fatalf("Failed to commit transaction: %v", commitErr)
			}
		}
	}()
	if _, err := tx.Exec(query, info.Name, info.Bio, info.Image, info.ID); err != nil {
		return err
	}
	return nil
}
