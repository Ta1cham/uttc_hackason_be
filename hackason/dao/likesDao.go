package dao

import (
	"database/sql"
	"log"
	"uttc_hackason_be/model"
)

type LikesDao struct {
	DB *sql.DB
}

func NewLikesDao(db *sql.DB) *LikesDao {
	return &LikesDao{DB: db}
}

func (dao *LikesDao) ToggleLikes(likes *model.LikeInfoPost) (bool, int, error) {
	db := dao.DB
	tx, err := db.Begin()
	if err != nil {
		return false, 0, err
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

	var countRecord int
	queryRecord := "SELECT COUNT(*) FROM likes WHERE uid=? AND tweet_id=?"
	if err := tx.QueryRow(queryRecord, likes.Uid, likes.TweetId).Scan(&countRecord); err != nil {
		return false, 0, err
	}

	var countLike int
	var isLike bool
	queryLikes := "SELECT COUNT(*) FROM likes WHERE tweet_id=?"
	if err := tx.QueryRow(queryLikes, likes.TweetId).Scan(&countLike); err != nil {
		return false, 0, err
	}

	// いいね済みなら取り消し　初ならいいね
	if countRecord > 0 {
		deleteQuery := "DELETE FROM likes WHERE uid=? AND tweet_id=?"
		if _, err := tx.Exec(deleteQuery, likes.Uid, likes.TweetId); err != nil {
			return false, 0, err
		}
		isLike = false
		countLike -= 1
	} else {
		insertQuery := "INSERT INTO likes (uid, tweet_id) VALUES (?, ?)"
		if _, err := tx.Exec(insertQuery, likes.Uid, likes.TweetId); err != nil {
			return false, 0, err
		}
		isLike = true
		countLike += 1
	}
	return isLike, countLike, nil
}
