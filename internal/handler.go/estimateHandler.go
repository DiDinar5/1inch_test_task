package handler

import (
	"net/http"

	"github.com/DiDinar5/1inch_test_task/domain"
	"github.com/labstack/echo/v4"
)

func (h *Handler) EstimateHandler(c echo.Context) error {
	var req domain.EstimateRequest
	//bind
	//validation middleware

	response, err := h.usecase.Estimate(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:       "Estimation failed",
			Code:        http.StatusInternalServerError,
			Description: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}
