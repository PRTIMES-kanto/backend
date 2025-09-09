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
タイトルはプレスリリースで最も重要な要素です。読者やメディアの関心を引き、ニュース価値（メディアフック）を前半に盛り込むことを重視してください。
以下の文章について、タイトル・リード文・本文それぞれに対して
1. "good": 現状で良いポイントを具体的に一文〜二文で褒めてください。タイトルでは読者やメディアの興味を引く要素があるかも評価してください。
2. "improvement": 現状どこが良くないかを具体的に指摘してください。特に、タイトル、リード文では、メディアフック（9つの要素:時流／季節性、画像／映像、逆説／対立、地域性、話題性、社会性／公益性、新規性／独自性、最上級／希少性、意外性）でどの要素が不足しているかを具体的に示してください。
3. "suggestion": 「次に何をすべきか」を具体的に指示し、その上でユーザーがすぐに書き直せる修正文例を1つ提示してください。
   - 必ず「次に〇〇してみましょう。その例として以下のように書き直せます。」の形式で書くこと。

必ず以下のJSONフォーマットのみを返してください。余計な文章・改行・説明は禁止です。

JSONフォーマット例:
{
  "title": {"good": "...", "improvement": "...", "suggestion": "..."},
  "lead":  {"good": "...", "improvement": "...", "suggestion": "..."},
  "body":  {"good": "...", "improvement": "...", "suggestion": "..."}
}

良いプレスリリースの定義を以下に定めます。これらの定義を基準にしてください:

- タイトル: 50-70文字程度で内容を要約し、5W2Hを盛り込みます。前半にはニュース性・話題性・数字・意外性などのメディアフックを含め、一目見て読者が『これ見たい!』と思うかどうかも評価してください。メディアフックとして以下の9つを意識してください:  
  1. 時流／季節性  
  2. 画像／映像  
  3. 逆説／対立  
  4. 地域性  
  5. 話題性  
  6. 社会性／公益性  
  7. 新規性／独自性  
  8. 最上級／希少性  
  9. 意外性  

- リード文: 250-300文字程度でまとめ、発信したい内容の5W2Hを盛り込み、プレスリリース全体が理解できる内容にしてください。加えて、メディアフック9つの要素が含まれているかも評価対象としてください。

- 本文: ① 誰に何のために伝えたい情報か、② なぜ自社がその活動を行うのか、③ 事実ベースで共感を呼ぶ内容、④ 理想の文書構成(起→承→転→展)。  
さらに、プレスリリースを通じて読者の共感を得るには、「感情」に訴えかけることが不可欠です。「喜び」「悲しみ」「恐怖・不安」「嫌悪」「驚き」「怒り」の6つの感情を引き出せるかを意識すると、メッセージがより深く伝わります。
- 喜び：嬉しい、楽しい、幸せな気持ち
- 悲しみ：共感ややさしさを誘う感情
- 恐怖・不安：問題提起とその解決策提示に活用
- 嫌悪：社会問題への関心喚起に
- 驚き：意外性やインパクトを与える
- 怒り：不満や理不尽さに共感し、行動を促す
プレスリリース作成時には、情報の読み手がどんな感情を抱くかを想定し、その感情に寄り添い、あるいは解消する立場で評価してください。例えば、読み手が「恐怖・不安」を感じる内容であればその対策を打ち出す、といったことです。

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
