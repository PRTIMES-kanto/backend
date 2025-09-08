package usecase

import (
	"prtimes/entity"
	"prtimes/external"
)

type ReviewUsecaseInterface interface {
	AnalyzeContent(title, lead, body string) (*entity.ReviewResult, error)
}

type ReviewUsecase struct {
	AIClient external.AIClientInterface
}

func NewReviewUsecase(aiClient external.AIClientInterface) ReviewUsecaseInterface {
	return &ReviewUsecase{
		AIClient: aiClient,
	}
}

func (u *ReviewUsecase) AnalyzeContent(title, lead, body string) (*entity.ReviewResult, error) {
	return u.AIClient.Analyze(title, lead, body)
}
