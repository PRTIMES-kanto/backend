package entity

type FieldReview struct {
	Good string `json:"good"`
	Improvement string `json:"improvement"`
	Suggestion  string `json:"suggestion"`
}

type ImageReview struct {
	URL         string `json:"url"`
    Good        string `json:"good"`
    Improvement string `json:"improvement"`
    Suggestion  string `json:"suggestion"`
}

// プレスリリース全体のレビュー結果
type ReviewResult struct {
	Title FieldReview `json:"title"`
	Lead  FieldReview `json:"lead"`
	Body  FieldReview `json:"body"`
	Image ImageReview `json:"image"`
}
