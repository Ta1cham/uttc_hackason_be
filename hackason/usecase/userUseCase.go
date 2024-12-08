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
	if data.Name == "" || utf8.RuneCountInString(data.Name) > 30 {
		return true
	}
	return false
}

func editInvalid(data *model.EditInfoForHTTPPOST) bool {
	if data.Name == "" || utf8.RuneCountInString(data.Name) > 30 {
		log.Printf(data.Name)
		return true
	}
	if bio := utf8.RuneCountInString(data.Bio); bio > 140 {
		log.Printf("editInvalid bio: %v", bio)
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

func (uc *UserUseCase) LoginUser(id string) (model.UserResForHTTPGET, error) {
	user, err := uc.UserDao.SearchUser(id)
	if err != nil {
		log.Printf("search user %s error: %v", id, err)
		return user, err
	}
	return user, nil
}

func (uc *UserUseCase) EditProfile(info *model.EditInfoForHTTPPOST) error {
	if editInvalid(info) {
		return fmt.Errorf("invalid user")
	}
	return uc.UserDao.EditProfile(info)
}
