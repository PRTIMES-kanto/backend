package main

import (
	// "log"
	// "os"

	"github.com/labstack/echo/v4"
	"prtimes/controller"
	"prtimes/external"
	"prtimes/usecase"
)

func main() {
	// Echo インスタンス作成
	e := echo.New()

	apiKey := ""

	// AIクライアント、ユースケース、コントローラの初期化
	aiClient := external.NewOpenAIClient(apiKey)
	reviewUsecase := usecase.NewReviewUsecase(aiClient)
	reviewController := controller.NewReviewController(reviewUsecase)

	// Ping用の確認ルート
	// e.GET("/ping", func(c echo.Context) error {
	// 	return c.JSON(200, map[string]string{"message": "pong"})
	// })

	// プレスリリースレビュー用エンドポイント
	e.POST("/review", reviewController.Review)

	// サーバ起動
	e.Logger.Fatal(e.Start(":8080"))
}
