package usecase

import (
	"errors"
	"github.com/oklog/ulid/v2"
	"log"
	"math/rand"
	"time"
	"unicode/utf8"
	"uttc_hackason_be/dao"
	"uttc_hackason_be/model"
)

func isInvalid(data model.UserInfoForHTTPPOST) bool {
	if data.Name == "" || utf8.RuneCountInString(data.Name) > 50 || data.Age < 20 || data.Age > 80 {
		return true
	}
	return false
}

func RegisterUser(dao *dao.UserDao) (ulid.ULID, error) {
	data := model.PostData
	if isInvalid(data) {
		log.Printf("")
		return ulid.ULID{}, errors.New("invalid data")
	}
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	if err := dao.RegisterUser(id.String()); err != nil {
		return ulid.ULID{}, err
	}
	return id, nil
}
