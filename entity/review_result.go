package entity

type FieldReview struct {
	Improvement string `json:"improvement"`
	Suggestion  string `json:"suggestion"`
}

// プレスリリース全体のレビュー結果
type ReviewResult struct {
	Title FieldReview `json:"title"`
	Lead  FieldReview `json:"lead"`
	Body  FieldReview `json:"body"`
}
