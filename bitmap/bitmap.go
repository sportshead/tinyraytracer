package bitmap

import (
	"fmt"
	"math"

	. "github.com/sportshead/tinyraytracer/vectors"
)

type Bitmap struct {
	Data   []byte
	Width  int
	Height int
}

func NewBitmap(width, height int) *Bitmap {
	return &Bitmap{
		Data:   make([]byte, width*height*4),
		Width:  width,
		Height: height,
	}
}

func (b *Bitmap) Set(x, y, index int, value byte) error {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height || index < 0 || index >= 4 {
		return fmt.Errorf("out of range: (%d, %d, %d)", x, y, index)
	}
	b.Data[(y*b.Width+x)*4+index] = value
	return nil
}

func (b *Bitmap) Get(x, y, index int) (byte, error) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height || index < 0 || index >= 4 {
		return 0, fmt.Errorf("out of range: (%d, %d, %d)", x, y, index)
	}
	return b.Data[(y*b.Width+x)*4+index], nil
}

func (b *Bitmap) SetPixel(x, y int, color Vec3f) error {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return fmt.Errorf("out of range: (%d, %d)", x, y)
	}
	for i := 0; i < 3; i++ {
		b.Set(x, y, i, floatToByte(color[i]))
	}
	b.Set(x, y, 3, 0xFF)
	return nil
}

func (b *Bitmap) GetPixel(x, y int) (Vec3f, error) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return Vec3f{}, fmt.Errorf("out of range: (%d, %d)", x, y)
	}
	var color Vec3f
	for i := 0; i < 3; i++ {
		value, err := b.Get(x, y, i)
		if err != nil {
			return Vec3f{}, err
		}
		color[i] = float64(value) / 255
	}
	return color, nil
}

func floatToByte(f float64) byte {
	return byte(math.Max(0, math.Min(255, f*255+0.5)))
}
