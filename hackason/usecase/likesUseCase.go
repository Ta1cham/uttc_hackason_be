package usecase

import (
	"database/sql"
	"uttc_hackason_be/dao"
	"uttc_hackason_be/model"
)

type LikesUseCase struct {
	LikesDao *dao.LikesDao
}

func NewLikesUseCase(db *sql.DB) *LikesUseCase {
	likesDao := dao.NewLikesDao(db)
	return &LikesUseCase{LikesDao: likesDao}
}

func (lc *LikesUseCase) ToggleLikes(likes *model.LikeInfoPost) (bool, int, error) {
	isLike, countLikes, err := lc.LikesDao.ToggleLikes(likes)
	if err != nil {
		return false, 0, err
	}
	return isLike, countLikes, nil
}
