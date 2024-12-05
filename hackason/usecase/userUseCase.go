package usecase

import (
	"database/sql"
	"fmt"
	"log"
	"unicode/utf8"
	"uttc_hackason_be/dao"
	"uttc_hackason_be/model"
)

type UserUseCase struct {
	UserDao *dao.UserDao
}

func NewUserUseCase(db *sql.DB) *UserUseCase {
	userDao := dao.NewUserDao(db)
	return &UserUseCase{UserDao: userDao}
}

func isInvalid(data *model.UserInfoForHTTPPOST) bool {
	if data.Name == "" || utf8.RuneCountInString(data.Name) > 50 || data.Age < 20 || data.Age > 80 {
		return true
	}
	return false
}

func (uc *UserUseCase) RegisterUser(user *model.UserInfoForHTTPPOST) error {
	if isInvalid(user) {
		return fmt.Errorf("invalid user")
	}
	return uc.UserDao.RegisterUser(user)
}

func (uc *UserUseCase) LoginUser(id string) ([]model.UserResForHTTPGET, error) {
	users, err := uc.UserDao.SearchUser(id)
	if err != nil {
		log.Printf("search user %s error: %v", id, err)
		return nil, err
	}
	return users, nil
}
