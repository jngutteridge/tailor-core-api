package core

type Text struct {
	Slug string `json:"slug"`
	Text string `json:"text"`
}

type TextResponse struct {
	Text []Text `json:"text"`
}