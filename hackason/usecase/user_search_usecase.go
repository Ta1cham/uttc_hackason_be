package usecase

import (
	"log"
	"uttc_hackason_be/dao"
	"uttc_hackason_be/model"
)

func SearchUserUsecase(name string, dao *dao.UserDao) ([]model.UserResForHTTPGET, error) {
	users, err := dao.SearchUser(name)
	if err != nil {
		log.Printf("search user %s error: %v", name, err)
		return nil, err
	}
	return users, nil
}
