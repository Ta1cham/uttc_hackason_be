package model

type TweetInfoForHTTPPOST struct {
	Uid     string `json:"uid"`
	Pid     string `json:"pid"`
	Content string `json:"content"`
	Imurl   string `json:"imurl"`
}

type TweetInfoForHTTPGET struct {
	Id       string  `json:"id"`
	Uid      string  `json:"uid"`
	Content  string  `json:"content"`
	Imurl    *string `json:"imurl"`
	PostedAt string  `json:"posted_at"`
	Uname    string  `json:"uname"`
	Likes    int     `json:"likes"`
	IsLike   bool    `json:"is_like"`
	Reps     int     `json:"reps"`
	Uimage   string  `json:"uimage"`
}
