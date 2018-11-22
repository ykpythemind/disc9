package disc9

import (
	"errors"
	"image"
	"io"
	"os"

	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

const maxDiscCount = 9

type Container struct {
	Discs    []Disc
	DiscSize int
	X        int
	Y        int
}

type Disc struct {
	Image image.Image
	UUID  string
	size  int
}

func (c *Container) ColorModel() color.Model {
	return c.Discs[0].Image.ColorModel()
}

func (c *Container) Bounds() image.Rectangle {
	min := image.Point{0, 0}
	max := image.Point{c.DiscSize - 1, c.DiscSize - 1}
	return image.Rectangle{
		Min: min,
		Max: max,
	}
}

func (c *Container) SaveImage(fileName string) error {
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	opts := &jpeg.Options{Quality: 100}

	return jpeg.Encode(out, c, opts)
}

type Pos struct {
	index int
	x     int
	y     int
}

func (c *Container) At(x, y int) color.Color {
	pos := detectIndexOfDisc(x, y, c.RectSize(), c.X, c.Y)
	disc := c.Discs[pos.index]
	discPx := x - (pos.x * c.RectSize())
	discPy := y - (pos.y * c.RectSize())

	return disc.Image.At(discPx, discPy)
}

func (c *Container) RectSize() int {
	return c.DiscSize / c.X
}

func detectIndexOfDisc(px, py, rectSize, xDisc, yDisc int) Pos {
	for x := 0; x < xDisc; x++ {
		for y := 0; y < yDisc; y++ {
			if (x*rectSize) <= px && px < (x+1)*rectSize {
				if (y*rectSize) <= py && py < (y+1)*rectSize {
					return Pos{detect(x, y), x, y} // FIXME
				}
			}
		}
	}
	return Pos{0, 0, 0}
}

func detect(x, y int) int {
	// FIXME
	if x == 0 {
		switch y {
		case 0:
			return 0
		case 1:
			return 3
		case 2:
			return 6
		}
	}
	if x == 1 {
		switch y {
		case 0:
			return 1
		case 1:
			return 4
		case 2:
			return 7
		}
	}
	switch y {
	case 0:
		return 2
	case 1:
		return 5
	}
	return 8
}

func NewContainer(rs []io.Reader, size, x, y int) (*Container, error) {
	if len(rs) != maxDiscCount {
		return nil, errors.New("io Must have 9")
	}

	discs := make([]Disc, 9)
	for i, r := range rs {
		d, err := NewDisc(r, size/x)
		if err != nil {
			return nil, err
		}
		d.Resize()
		discs[i] = *d
	}

	container := &Container{
		Discs:    discs,
		DiscSize: size,
		X:        x,
		Y:        y,
	}
	return container, nil
}

func (d *Disc) Resize() {
	w := uint(d.size)
	d.Image = resize.Resize(w, w, d.Image, resize.Lanczos3)
}

func NewDisc(r io.Reader, size int) (*Disc, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	disc := &Disc{
		Image: img,
		UUID:  uuid.New().String(),
		size:  size,
	}

	return disc, nil
}
