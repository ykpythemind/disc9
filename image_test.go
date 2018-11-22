package disc9

import (
	"io"
	"os"
	"strconv"
	"testing"
)

func TestSaveImage(t *testing.T) {
	var rs [9]io.Reader
	for i := range rs {
		s := strconv.Itoa(i + 1)
		r, err := os.Open("./testdata/" + s + ".jpg")
		defer r.Close()
		if err != nil {
			panic(err)
		}
		rs[i] = r
	}

	con, err := NewContainer(rs[:], 500, 3, 3)
	if err != nil {
		t.Error(err)
	}

	err = con.SaveImage("test.jpg")
	if err != nil {
		t.Error(err)
	}
}
