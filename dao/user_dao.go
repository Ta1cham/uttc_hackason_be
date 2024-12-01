package dao

import (
	"database/sql"
	"uttc_hackason_be/model"
)

type UserDao struct {
	DB *sql.DB
}

func CreateDao(db *sql.DB) *UserDao {
	dao := &UserDao{DB: db}
	return dao
}

func (dao *UserDao) RegisterUser(id string) error {
	db := dao.DB
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	data := model.PostData
	if _, err := tx.Exec("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", id, data.Name, data.Age); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (dao *UserDao) SearchUser(name string) ([]model.UserResForHTTPGET, error) {
	db := dao.DB
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
}
