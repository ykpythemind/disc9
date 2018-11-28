package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo"
)

func TestUploadHandler(t *testing.T) {
	t.Skip()
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	file, err := os.Open("testdata/1.jpg")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	for i := 0; i < 9; i++ {
		// name := "files[" + string(i) + "]"
		// part := make(textproto.MIMEHeader)
		// part.Set("Content-Type", "image/jpeg")
		// part.Set("Content-Disposition", `form-data; name="files"; filename="1.jpg"`) // works
		// part.Set("Content-Disposition", `form-data; name="`+name+`"; filename="1.jpg"`)
		// writer, err := w.CreatePart(part)
		writer, err := w.CreateFormFile("files", "1.jpg")
		if err != nil {
			t.Error(err)
		}
		if _, err := io.Copy(writer, file); err != nil {
			t.Error(err)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	w.Close()

	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	if err := uploadHandler(c); err != nil {
		t.Errorf("error! %s", err)
	}
}
