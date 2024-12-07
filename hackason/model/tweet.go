package model

type TweetInfoForHTTPPOST struct {
	Uid     string `json:"uid"`
	Content string `json:"content"`
	Imurl   string `json:"imurl"`
}

type TweetInfoForHTTPGET struct {
	Id       string  `json:"id"`
	Uid      string  `json:"uid"`
	Content  string  `json:"content"`
	Imurl    *string `json:"imurl"`
	Likes    int     `json:"likes"`
	PostedAt string  `json:"posted_at"`
	Uname    string  `json:"uname"`
}
