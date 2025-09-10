package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	"prtimes/entity"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

const s3Bucket = "gazou-hozon"

type AIClientInterface interface {
	Analyze(title, lead, body, s3ImageURL string) (*entity.ReviewResult, error)
	UploadImageToS3(imageURL string) (string, error)
}

type OpenAIClient struct {
	client openai.Client
	s3Client *s3.Client
	bucket string
}

func NewOpenAIClient(apiKey string) AIClientInterface {
	// AWS S3クライアント初期化
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-1"),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	return &OpenAIClient{
		client: openai.NewClient(option.WithAPIKey(apiKey)),
		s3Client: s3Client,
		bucket:   s3Bucket,
	}
}

// URLから画像をダウンロードしてS3にアップロード
func (o *OpenAIClient) UploadImageToS3(imageURL string) (string, error) {
	// 1. 画像をダウンロード
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// 2. ユニークなファイル名を生成
	ext := path.Ext(imageURL)
	if ext == "" {
		ext = ".jpg" // デフォルト拡張子
	}
	key := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 3. S3 にアップロード
	_, err = o.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(o.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
		ACL:    types.ObjectCannedACLPublicRead, // 公開URL
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// 4. S3 の公開 URL を返す
	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", o.bucket, key)
	fmt.Println("S3保存処理完了:", s3URL)
	return s3URL, nil
}



func (o *OpenAIClient) Analyze(title, lead, body, s3ImageURL string) (*entity.ReviewResult, error) {
    // fmt.Printf("DEBUG: title: %s, lead: %s, body: %s, imageURL: %s\n", 
    //     title, lead, body, s3ImageURL)
	prompt := fmt.Sprintf(`
あなたはプレスリリースのレビュアーです。
タイトルはプレスリリースで最も重要な要素です。読者やメディアの関心を引き、ニュース価値（メディアフック）を前半に盛り込むことを重視してください。
以下の文章について、タイトル・リード文・本文それぞれに対して
1. "good": 現状で良いポイントを具体的に一文〜二文で褒めてください。必ず**理由や根拠**を示すこと。「なぜ良いと感じたのか」を具体的に説明してください。
2. "improvement": 現状どこが良くないかを具体的に指摘してください。特に、タイトル、リード文では、メディアフック（9つの要素:時流／季節性、画像／映像、逆説／対立、地域性、話題性、社会性／公益性、新規性／独自性、最上級／希少性、意外性）でどの要素が不足しているかを**根拠付きで具体的に示すこと**。
3. "suggestion": 「次に何をすべきか」を具体的に指示し、その上でユーザーがすぐに書き直せる修正文例を1つ提示してください。必ず**理由や意図を添えること**。
   - 形式は「次に〇〇してみてはいかがでしょうか。その例として以下が考えられます。」としてください。

良いプレスリリースの定義を以下に定めます。これらの定義を基準にしてください:

- タイトル: 50-70文字程度で内容を要約し、5W2Hを盛り込みます。前半にはニュース性・話題性・数字・意外性などのメディアフックを含み、一目見て読者が『これ見たい!』と思うかどうかも評価してください。メディアフックとして以下の9つを意識してください:  
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
  また、感情に訴えかけることも評価してください。「喜び」「悲しみ」「恐怖・不安」「嫌悪」「驚き」「怒り」の6つの感情を引き出せるかを意識してください。

さらに、タイトル直後に掲載されるメイン画像1枚も評価してください。
- 評価対象は「タイトル直後に掲載されるメイン画像1枚のみ」です。
- メディアが使いやすいか、情報が伝わりやすいかを評価する際には、必ず**根拠を添えて説明してください**。
- 評価は以下の観点を含めてください:
  1. 十分なサイズ・解像度があるか
  2. 構図や視線誘導で情報が伝わりやすいか
  3. 商品・サービスやイベントの全体像が理解できるか
  4. 記事に掲載する際の使いやすさ

必ず以下のJSONフォーマットのみを返してください。余計な文章・改行・説明は禁止です。

JSONフォーマット例:
{
  "title": {"good": "...根拠付きで...", "improvement": "...根拠付きで...", "suggestion": "...根拠付きで次のアクション..."},
  "lead":  {"good": "...根拠付きで...", "improvement": "...根拠付きで...", "suggestion": "...根拠付きで次のアクション..."},
  "body":  {"good": "...根拠付きで...", "improvement": "...根拠付きで...", "suggestion": "...根拠付きで次のアクション..."},
  "image": {
      "url": "S3にアップロードした画像URL",
      "good": "...根拠付きで...",
      "improvement": "...根拠付きで...",
      "suggestion": "...根拠付きで次のアクション..."
  }
}

画像URL: %s
入力文章:
{
  "title": "%s",
  "lead": "%s",
  "body": "%s"
}

`, s3ImageURL, title, lead, body)

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
