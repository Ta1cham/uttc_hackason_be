package model

type UserResForHTTPGET struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserInfoForHTTPPOST struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var PostData UserInfoForHTTPPOST
