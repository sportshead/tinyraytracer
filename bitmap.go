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
