package dao

import (
	"database/sql"
	"log"
	"uttc_hackason_be/model"
)

type TweetDao struct {
	DB *sql.DB
}

func NewTweetDao(db *sql.DB) *TweetDao {
	dao := &TweetDao{DB: db}
	return dao
}

func (dao *TweetDao) MakeTweet(tweet *model.TweetInfoForHTTPPOST) error {
	db := dao.DB
	tx, err := db.Begin()
	if err != nil {
		return err
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

	if _, err := tx.Exec("INSERT INTO tweet (uid, content, image) VALUES (?, ?, ?)", tweet.Uid, tweet.Content, tweet.Imurl); err != nil {
		log.Printf("Failed to insert tweet: %v", err)
		return err
	}

	return nil
}

func (dao *TweetDao) GetTweet(pNum int, currentUser string) ([]model.TweetInfoForHTTPGET, error) {
	db := dao.DB

	// 一度に読み込むページ数
	const pageSize = 10

	query := `
        SELECT tweet.id, tweet.uid, tweet.content, tweet.image, tweet.posted_at,
               user.name AS uname
        FROM tweet
        JOIN user ON user.id = tweet.uid
        ORDER BY posted_at DESC
        LIMIT ? OFFSET ?;
    `

	rows, err := db.Query(query, pageSize, pageSize*pNum)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []model.TweetInfoForHTTPGET
	for rows.Next() {
		var tweet model.TweetInfoForHTTPGET
		if err := rows.Scan(&tweet.Id, &tweet.Uid, &tweet.Content, &tweet.Imurl, &tweet.PostedAt, &tweet.Uname); err != nil {
			return nil, err
		}

		var likesCount int
		likesCountQuery := `
			SELECT COUNT(*) 
			FROM likes 
			WHERE tweet_id = ?;
		`
		err := db.QueryRow(likesCountQuery, tweet.Id).Scan(&likesCount)
		if err != nil {
			log.Printf("Failed to get likes: %v", err)
			return nil, err
		}
		tweet.Likes = likesCount

		// 現在のユーザーがこのツイートにいいねしているかを確認
		var isLiked int
		isLikedQuery := `
			SELECT COUNT(*) 
			FROM likes 
			WHERE tweet_id = ? AND uid = ?;
		`
		err = db.QueryRow(isLikedQuery, tweet.Id, currentUser).Scan(&isLiked)
		if err != nil {
			log.Printf("Failed to get isLiked: %v", err)
			return nil, err
		}
		tweet.IsLike = isLiked > 0

		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tweets, nil
}
