package handler

import (
	repository "core/internal/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type IndicatorHandler struct {
	repo *repository.IndicatorRepository
}

func NewIndicatorHandler(repo *repository.IndicatorRepository) *IndicatorHandler {
	return &IndicatorHandler{repo: repo}
}

func (h *IndicatorHandler) GetIndicatorsHandler(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	page, _ := strconv.Atoi(c.QueryParam("page"))

	if limit <= 0 {
		limit = 10
	}

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit
	indicators, err := h.repo.GetIndicators(limit, offset)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Indicators not found"})
	}

	return c.JSON(http.StatusOK, indicators)
}

func (h *IndicatorHandler) GetIndicatorByIdHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid indicator ID"})
	}

	indicator, err := h.repo.GetIndicatorById(id)
	if err != nil {
		if err.Error() == "indicator not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Indicator not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, indicator)
}
