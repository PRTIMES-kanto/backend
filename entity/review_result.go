package entity

type ReviewResult struct {
	TitleFeedback string `json:"title_feedback"`
	LeadFeedback  string `json:"lead_feedback"`
	BodyFeedback  string `json:"body_feedback"`
}