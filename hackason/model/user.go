package model

type UserResForHTTPGET struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Bio   string `json:"bio"`
	Image string `json:"image"`
}

type UserInfoForHTTPPOST struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EditInfoForHTTPPOST struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Bio   string `json:"bio"`
}
