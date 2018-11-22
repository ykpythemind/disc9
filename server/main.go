package main

import (
	"io"
	"net/http"

	disc9 "github.com/ykpythemind/9disc"

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
	// Read form fields
	// name := c.FormValue("name")
	// email := c.FormValue("email")

	//------------
	// Read files
	//------------

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	var readers [9]io.Reader
	var count = 0

	for i, file := range files {
		count++
		if count > 9 {
			break
		}
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		readers[i] = src
	}

	if count != 9 {
		return c.String(http.StatusUnprocessableEntity, "9 files are required")
	}

	con, err := disc9.NewContainer(readers[:], 500, 3, 3)
	if err != nil {
		return err
	}

	return con.ToJpeg(c.Response())
}
