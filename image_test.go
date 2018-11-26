package disc9

import (
	"io"
	"io/ioutil"
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

func BenchmarkToJpeg(b *testing.B) {
	readers, closes := prepareReaders()
	defer closes()

	con, err := NewContainer(readers[:], 500)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		err = con.ToJpeg(ioutil.Discard)
		if err != nil {
			panic(err)
		}
	}
}

func prepareReaders() ([]io.Reader, func()) {
	var readers [9]io.Reader
	closes := [9]func() error{}
	for i := range readers {
		s := strconv.Itoa(i + 1)
		r, err := os.Open("./testdata/" + s + ".jpg")
		closes[i] = r.Close
		if err != nil {
			panic(err)
		}
		readers[i] = r
	}

	return readers[:], func() {
		for _, c := range closes {
			c()
		}
	}
}

func TestSaveImage(t *testing.T) {
	readers, closes := prepareReaders()
	defer closes()

	con, err := NewContainer(readers[:], 500)
	if err != nil {
		t.Error(err)
	}

	err = con.SaveImage("test.jpg")
	if err != nil {
		t.Error(err)
	}
}
