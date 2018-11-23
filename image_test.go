package disc9

import (
	"io"
	"os"
	"strconv"
	"testing"
)

func TestPosition(t *testing.T) {
	var tests = []struct {
		expected int
		given    position
	}{
		{0, position{0, 0}},
		{3, position{0, 1}},
		{4, position{1, 1}},
		{8, position{2, 2}},
	}
	for i, tt := range tests {
		tt := tt
		t.Run("test "+string(i), func(t *testing.T) {
			actual := (tt.given.index())
			if actual != tt.expected {
				t.Errorf("(%s): expected %d, actual %d", tt.given.String(), tt.expected, actual)
			}

		})
	}
}

func TestSaveImage(t *testing.T) {
	var readers [9]io.Reader
	for i := range readers {
		s := strconv.Itoa(i + 1)
		r, err := os.Open("./testdata/" + s + ".jpg")
		defer r.Close()
		if err != nil {
			panic(err)
		}
		readers[i] = r
	}

	con, err := NewContainer(readers[:], 500, 3, 3)
	if err != nil {
		t.Error(err)
	}

	err = con.SaveImage("test.jpg")
	if err != nil {
		t.Error(err)
	}
}
