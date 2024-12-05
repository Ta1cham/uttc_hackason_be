package dao

import (
	"database/sql"
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
	if _, err := tx.Exec("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", user.ID, user.Name, user.Age); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (dao *UserDao) SearchUser(id string) ([]model.UserResForHTTPGET, error) {
	db := dao.DB
	tx, err := db.Begin()
	if err != nil {
		log.Printf("")
		return nil, err
	}
	rows, err := tx.Query("SELECT id, name, age FROM user WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Printf("")
		return nil, err
	}
	users := make([]model.UserResForHTTPGET, 0)
	for rows.Next() {
		var u model.UserResForHTTPGET
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			log.Printf("fail: rows.Scan, %v\n", err)

			if err := rows.Close(); err != nil {
				log.Printf("fail: rows.Close(), %v\n", err)
			}
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Close(); err != nil {
		log.Printf("rows.Close(), %v\n", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		log.Printf("")
		return nil, err
	}
	return users, nil
}
