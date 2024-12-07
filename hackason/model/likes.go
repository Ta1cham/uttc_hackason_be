package model

type LikeInfoPost struct {
	TweetId string `json:"tweet_id"`
	Uid     string `json:"uid"`
}

type LikeResponse struct {
	IsLike    bool `json:"is_like"`
	CountLike int  `json:"count_like"`
}
