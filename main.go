package main

import (

	"github.com/labstack/echo/v4"
	"prtimes/controller"
	"prtimes/external"
	"prtimes/usecase"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo インスタンス作成
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

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
