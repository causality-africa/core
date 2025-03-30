package api

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

var (
	pageSize    = 25
	maxPageSize = 50
)

type pagination struct {
	Page int
	Size int
}

func getPaginationParams(c echo.Context) *pagination {
	var params pagination
	var err error

	params.Page, err = strconv.Atoi(c.QueryParam("page"))
	if err != nil || params.Page < 1 {
		params.Page = 1
	}

	params.Size, err = strconv.Atoi(c.QueryParam("size"))
	if err != nil || params.Size < 1 || params.Size > maxPageSize {
		params.Size = pageSize
	}

	return &params
}
