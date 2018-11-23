package main

import (
	"io"
	"net/http"

	"github.com/ykpythemind/disc9"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start("localhost:1323"))
}

func upload(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	if len(files) != 9 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "9 files are required")
	}

	var readers [9]io.Reader

	for i, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		readers[i] = src
	}

	container, err := disc9.NewContainer(readers[:], 500, 3, 3)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	return container.ToJpeg(c.Response())
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func customErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	e := errorResponse{
		Code:    code,
		Message: err.Error(),
	}

	c.JSON(code, e)
}
