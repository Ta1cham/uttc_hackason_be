package model

type UserResForHTTPGET struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserInfoForHTTPPOST struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
