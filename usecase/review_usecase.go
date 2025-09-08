package usecase

import (
	"prtimes/entity"
	"prtimes/external"
)

type ReviewUsecaseInterface interface {
	AnalyzeContent(content string) (*entity.ReviewResult, error)
}

type ReviewUsecase struct {
	AIClient external.AIClientInterface
}

func NewReviewUsecase(aiClient external.AIClientInterface) ReviewUsecaseInterface {
	return &ReviewUsecase{
		AIClient: aiClient,
	}
}

func (u *ReviewUsecase) AnalyzeContent(content string) (*entity.ReviewResult, error) {
	return u.AIClient.Analyze(content)
}