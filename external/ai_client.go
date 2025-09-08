package external

import "prtimes/entity"

type AIClientInterface interface {
	// これってなんでentityにポインタつけているんだっけ？
	Analyze(content string) (*entity.ReviewResult, error)
}

type MockAIClient struct {}

func NewMockAIClient() AIClientInterface {
	return &MockAIClient{}
}

func (m *MockAIClient) Analyze(content string) (*entity.ReviewResult, error) {
	// 固定レスポンスを返す
	return &entity.ReviewResult {
		TitleFeedback: "タイトルをもう少し具体的にすると良いです",
		LeadFeedback:  "リード文に5W2Hを含めるとさらに良いです",
		BodyFeedback:  "本文に背景情報を加えると説得力が増します",
	}, nil
}