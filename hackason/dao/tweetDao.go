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
	var query string
	var args []interface{}

	// pidがNULLかどうかでクエリと引数を分ける
	if tweet.Pid == "" {
		query = "INSERT INTO tweet (uid, content, image) VALUES (?, ?, ?)"
		args = []interface{}{tweet.Uid, tweet.Content, tweet.Imurl}
	} else {
		query = "INSERT INTO tweet (uid, content, image, pid) VALUES (?, ?, ?, ?)"
		args = []interface{}{tweet.Uid, tweet.Content, tweet.Imurl, tweet.Pid}
	}

	if _, err := tx.Exec(query, args...); err != nil {
		log.Printf("Failed to insert tweet: %v", err)
		return err
	}

	return nil
}

func (dao *TweetDao) GetTweet(pNum int, currentUser string, pid string) ([]model.TweetInfoForHTTPGET, error) {
	db := dao.DB

	// 一度に読み込むページ数
	const pageSize = 10

	baseQuery := `
        SELECT tweet.id, tweet.uid, tweet.content, tweet.image, tweet.posted_at,
               user.name AS uname
        FROM tweet
        JOIN user ON user.id = tweet.uid
    `

	var query string
	var args []interface{}
	if pid != "" {
		query = baseQuery + "WHERE tweet.pid = ? ORDER BY posted_at DESC LIMIT ? OFFSET ?;"
		args = append(args, pid, pageSize, pageSize*pNum)
	} else {
		query = baseQuery + "WHERE tweet.pid IS NULL ORDER BY posted_at DESC LIMIT ? OFFSET ?;"
		args = append(args, pageSize, pageSize*pNum)
	}

	rows, err := db.Query(query, args...)
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

		var reps int
		repsQuery := `
			SELECT COUNT(*) 
			FROM tweet
			WHERE pid = ?
		`
		err = db.QueryRow(repsQuery, tweet.Id).Scan(&reps)
		if err != nil {
			log.Printf("Failed to get reps: %v", err)
			return nil, err
		}
		tweet.Reps = reps

		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tweets, nil
}

func (dao *TweetDao) GetTweetById(id string, currentUser string) (model.TweetInfoForHTTPGET, error) {
	db := dao.DB
	var tweet model.TweetInfoForHTTPGET
	tx, err := db.Begin()
	if err != nil {
		return model.TweetInfoForHTTPGET{}, err
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

	query := `
		SELECT tweet.id, tweet.uid, tweet.content, tweet.image, tweet.posted_at,
               user.name AS uname
        FROM tweet
        JOIN user ON user.id = tweet.uid
		WHERE tweet.id = ?;`

	err = tx.QueryRow(query, id).Scan(&tweet.Id, &tweet.Uid, &tweet.Content, &tweet.Imurl, &tweet.PostedAt, &tweet.Uname)
	if err != nil {
		return model.TweetInfoForHTTPGET{}, err
	}

	var likesCount int
	likesCountQuery := `
			SELECT COUNT(*) 
			FROM likes 
			WHERE tweet_id = ?;
		`
	err = db.QueryRow(likesCountQuery, tweet.Id).Scan(&likesCount)
	if err != nil {
		log.Printf("Failed to get likes: %v", err)
		return model.TweetInfoForHTTPGET{}, err
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
		return model.TweetInfoForHTTPGET{}, err
	}
	tweet.IsLike = isLiked > 0

	var reps int
	repsQuery := `
			SELECT COUNT(*) 
			FROM tweet
			WHERE pid = ?
		`
	err = tx.QueryRow(repsQuery, tweet.Id).Scan(&reps)
	if err != nil {
		log.Printf("Failed to get reps: %v", err)
		return model.TweetInfoForHTTPGET{}, err
	}
	tweet.Reps = reps

	return tweet, nil
}
