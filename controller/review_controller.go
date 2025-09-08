package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"prtimes/usecase"
)

type ReviewController struct {
	Usecase usecase.ReviewUsecaseInterface
}

func NewReviewController(u usecase.ReviewUsecaseInterface) *ReviewController {
	return &ReviewController{
		Usecase: u,
	}
}

func (rc *ReviewController) Review(c echo.Context) error {
	var req struct {
		Title string `json:"title"`
		Lead  string `json:"lead"`
		Body  string `json:"body"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	result, err := rc.Usecase.AnalyzeContent(req.Title, req.Lead, req.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "analysis failed", "detail": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
