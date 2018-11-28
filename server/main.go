package main

import (
	"fmt"
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
	e.POST("/upload", uploadHandler)

	e.Logger.Fatal(e.Start("localhost:1323"))
}

func uploadHandler(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	if len(files) != 9 {
		s := fmt.Sprintf("9 files are required, but %d files given", len(files))
		return echo.NewHTTPError(http.StatusUnprocessableEntity, s)
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

	container, err := disc9.NewContainer(readers[:], 500)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
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
