package disc9

import (
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"sync"

	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
)

const maxDiscCount = 9

type Container struct {
	Discs    []*disc
	DiscSize int
	XCount   int
	YCount   int
}

type disc struct {
	Image image.Image
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

func (c *Container) ToJpeg(out io.Writer) error {
	opts := &jpeg.Options{Quality: 100}
	return jpeg.Encode(out, c, opts)
}

func (c *Container) SaveImage(fileName string) error {
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	return c.ToJpeg(out)
}

type position struct {
	x int
	y int
}

func (p *position) String() string {
	return fmt.Sprintf("x = %+v, y = %+v", p.x, p.y)
}

var posMap = [][]int{
	[]int{0, 1, 2},
	[]int{3, 4, 5},
	[]int{6, 7, 8},
}

func (p *position) index() int {
	return posMap[p.y][p.x]
}

func (c *Container) At(x, y int) color.Color {
	pos := c.detectPositionFromPoint(x, y)
	disc := c.Discs[pos.index()]
	discPx := x - (pos.x * c.rectSize())
	discPy := y - (pos.y * c.rectSize())

	return disc.Image.At(discPx, discPy)
}

func (c *Container) rectSize() int {
	return c.DiscSize / c.XCount
}

func (c *Container) detectPositionFromPoint(px, py int) position {
	xDisc := c.XCount
	yDisc := c.YCount
	rectSize := c.rectSize()

	for x := 0; x < xDisc; x++ {
		for y := 0; y < yDisc; y++ {
			if (x*rectSize) <= px && px < (x+1)*rectSize {
				if (y*rectSize) <= py && py < (y+1)*rectSize {
					return position{x, y}
				}
			}
		}
	}
	return position{0, 0}
}

// NewContainer is constructor of Container
//  size: output square size(px)
func NewContainer(readers []io.Reader, size int) (*Container, error) {
	// TODO: xとy方向の枚数を変更可能にする
	x := 3
	y := 3

	if len(readers) != maxDiscCount {
		return nil, errors.New("io Must have 9")
	}

	discs := make([]*disc, maxDiscCount)

	wg := sync.WaitGroup{}
	for i, r := range readers {
		d, err := newDisc(r, size/x)
		if err != nil {
			return nil, err
		}
		wg.Add(1)
		go func() {
			d.resize() // Resize each disc size here
			wg.Done()
		}()
		discs[i] = d
	}
	wg.Wait()

	container := &Container{
		Discs:    discs,
		DiscSize: size,
		XCount:   x,
		YCount:   y,
	}
	return container, nil
}

func (d *disc) resize() {
	w := uint(d.size)
	d.Image = resize.Resize(w, w, d.Image, resize.Lanczos3)
}

func newDisc(r io.Reader, size int) (*disc, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	disc := &disc{
		Image: img,
		size:  size,
	}

	return disc, nil
}
