package external

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"prtimes/entity"
)

type AIClientInterface interface {
	Analyze(title, lead, body string) (*entity.ReviewResult, error)
}

type OpenAIClient struct {
	client openai.Client
}

func NewOpenAIClient(apiKey string) AIClientInterface {
	return &OpenAIClient{
		client: openai.NewClient(option.WithAPIKey(apiKey)),
	}
}

func (o *OpenAIClient) Analyze(title, lead, body string) (*entity.ReviewResult, error) {
	prompt := fmt.Sprintf(`
あなたはプレスリリースのレビュアーです。
以下の文章について、タイトル・リード文・本文それぞれに対して
1. "improvement": 現状どこが良くないかを具体的に指摘
2. "suggestion": 「次に何をすべきか」を具体的に指示し、その上でユーザーがすぐに書き直せる修正文例を1つ提示してください。
   - 必ず「次に〇〇してみましょう。その例として以下のように書き直せます。」の形式で書くこと。

必ず以下のJSONフォーマットのみを返してください。余計な文章・改行・説明は禁止です。

JSONフォーマット例:
{
  "title": {"improvement": "...", "suggestion": "..."},
  "lead": {"improvement": "...", "suggestion": "..."},
  "body": {"improvement": "...", "suggestion": "..."}
}

良いプレスリリースの定義を以下に定めます。これらの定義を基準にしてください:
- タイトル: - 50-70文字程度で内容を要約できている、5W2Hから伝えるべき内容を盛り込む(Who,What,How muchという情報)、メディアフック=ニュースバリューを前方30文字程度に(意外性、話題性、3分という数字のインパクト、AIという注目度の高いキーワード)。
- リード文: 250-300文字程度でまとめる、発信したい内容の5W2Hを盛り込む、リード文だけでプレスリリース全体の内容が理解できるように。
- 本文: ① 誰に、何のために、伝えたい情報なのか、② なぜ、自社がその活動を行うのか、③事実をベースに、でも共感を呼ぶ内容に、④ 本文での理想の文書構成(起：企業活動の背景・きっかけ) → (承：社会背景や課題、統計データなどがあるとなお良し) → (転:取組む中での転換点があればぜひ記載、担当者コメントなどで盛り込めると共感性がUP) → (展：今後の展望。継続的なステークホルダーを増やす)

入力文章:
{
  "title": "%s",
  "lead": "%s",
  "body": "%s"
}
`, title, lead, body)

	resp, err := o.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4o,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	})
	if err != nil {
		return nil, err
	}

	aiText := resp.Choices[0].Message.Content

	// 余計な ```json ``` があれば削除
	aiText = strings.TrimSpace(aiText)
	aiText = strings.TrimPrefix(aiText, "```json")
	aiText = strings.TrimSuffix(aiText, "```")
	aiText = strings.TrimSpace(aiText)

	var result entity.ReviewResult
	if err := json.Unmarshal([]byte(aiText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w\nAI response: %s", err, aiText)
	}

	return &result, nil
}
