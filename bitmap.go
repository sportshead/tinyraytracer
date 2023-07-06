package main

import "errors"

type Bitmap struct {
	data   []byte
	width  int
	height int
}

func NewBitmap(width, height int) *Bitmap {
	return &Bitmap{
		data:   make([]byte, width*height*4),
		width:  width,
		height: height,
	}
}

func (b *Bitmap) Width() int {
	return b.width
}

func (b *Bitmap) Height() int {
	return b.height
}

func (b *Bitmap) Data() []byte {
	return b.data
}

func (b *Bitmap) Set(x, y, index int, value byte) error {
	if x < 0 || x >= b.width || y < 0 || y >= b.height || index < 0 || index >= 4 {
		return errors.New("out of range")
	}
	b.data[(y*b.width+x)*4+index] = value
	return nil
}

func (b *Bitmap) Get(x, y, index int) (byte, error) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height || index < 0 || index >= 4 {
		return 0, errors.New("out of range")
	}
	return b.data[(y*b.width+x)*4+index], nil
}

func (b *Bitmap) SetPixel(x, y int, color Vec3f) error {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return errors.New("out of range")
	}
	for i := 0; i < 3; i++ {
		b.Set(x, y, i, ftob(color[i]))
	}
	b.Set(x, y, 3, 0xFF)
	return nil
}

func (b *Bitmap) GetPixel(x, y int) (Vec3f, error) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return Vec3f{}, errors.New("out of range")
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
