package usecase

import (
	"database/sql"
	"errors"
	"unicode/utf8"
	"uttc_hackason_be/dao"
	"uttc_hackason_be/model"
)

type TweetUseCase struct {
	TweetDao *dao.TweetDao
}

func NewTweetUseCase(db *sql.DB) *TweetUseCase {
	tweetDao := dao.NewTweetDao(db)
	return &TweetUseCase{TweetDao: tweetDao}
}

// htmlタグのエスケープとか必要そうならやる
func checkTweetPosted(tweet *model.TweetInfoForHTTPPOST) error {
	if tweet.Content == "" || utf8.RuneCountInString(tweet.Content) > 140 {
		return errors.New("invalid tweet content")
	}
	return nil
}

func (tc *TweetUseCase) MakeTweet(tweet *model.TweetInfoForHTTPPOST) error {
	if err := checkTweetPosted(tweet); err != nil {
		return err
	}
	return tc.TweetDao.MakeTweet(tweet)
}

//func checkTweetGet(tweets []model.TweetInfoForHTTPGET) {
//	for _, tweet := range tweets {
//	// 削除されたユーザの投稿は非表示
//	tweet.
//	}
//	return
//}

func (tc *TweetUseCase) GetTweet(pNum int, currentUser string, pid string) ([]model.TweetInfoForHTTPGET, error) {
	tweets, err := tc.TweetDao.GetTweet(pNum, currentUser, pid)
	if err != nil {
		return nil, err
	}
	return tweets, nil
}

func (tc *TweetUseCase) GetTweetById(id string, currentUser string) (model.TweetInfoForHTTPGET, error) {
	tweet, err := tc.TweetDao.GetTweetById(id, currentUser)
	if err != nil {
		return model.TweetInfoForHTTPGET{}, err
	}
	return tweet, nil
}
