package usecase

import (
	"prtimes/entity"
	"prtimes/external"
)

type ReviewUsecaseInterface interface {
	AnalyzeContent(title, lead, body, mainImageURL string) (*entity.ReviewResult, error)
}

type ReviewUsecase struct {
	AIClient external.AIClientInterface
}

func NewReviewUsecase(aiClient external.AIClientInterface) ReviewUsecaseInterface {
	return &ReviewUsecase{
		AIClient: aiClient,
	}
}

func (u *ReviewUsecase) AnalyzeContent(title, lead, body, mainImageURL string) (*entity.ReviewResult, error) {
	// 画像をS3にアップロードしてURLを取得する
	s3URL, err := u.AIClient.UploadImageToS3(mainImageURL)
	if err != nil {
		return nil, err
	}

	return u.AIClient.Analyze(title, lead, body, s3URL)
}
