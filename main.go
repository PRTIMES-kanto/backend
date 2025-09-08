package main

import (
	"github.com/labstack/echo/v4"
	"prtimes/controller"
	"prtimes/external"
	"prtimes/usecase"
)

func main() {
	e := echo.New()

	aiClient := external.NewMockAIClient()
	reviewUsecase := usecase.NewReviewUsecase(aiClient)
	reviewController := controller.NewReviewController(reviewUsecase)

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "pong"})
	})

	e.POST("/review", reviewController.Review)

	e.Logger.Fatal(e.Start(":8080"))
}
