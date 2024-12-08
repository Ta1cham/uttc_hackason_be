package model

type NoteForPost struct {
	Pid  string `json:"pid"`
	Note string `json:"note"`
}

type ContentData struct {
	Pid     string `json:"pid"`
	Content string `json:"content"`
}

type Response struct {
	Candidates    []Candidate   `json:"Candidates"`
	UsageMetadata UsageMetadata `json:"UsageMetadata"`
}

type Candidate struct {
	Index        int     `json:"Index"`
	Content      Content `json:"Content"`
	FinishReason int     `json:"FinishReason"`
}

type Content struct {
	Role  string   `json:"Role"`
	Parts []string `json:"Parts"`
}

type UsageMetadata struct {
	PromptTokenCount     int `json:"PromptTokenCount"`
	CandidatesTokenCount int `json:"CandidatesTokenCount"`
	TotalTokenCount      int `json:"TotalTokenCount"`
}
