package icon

import (
	"bytes"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

var (
	colorRed    = color.RGBA{R: 227, G: 78, B: 73, A: 255}
	colorYellow = color.RGBA{R: 242, G: 211, B: 36, A: 255}
	colorGreen  = color.RGBA{R: 73, G: 204, B: 130, A: 255}
	colorWhite  = color.RGBA{R: 255, G: 255, B: 255, A: 255}
)

type AnimateIcon interface {
	Update() error
}

type TextIconImage struct {
	Text  *string
	Color *color.RGBA
}

func (textIcon *TextIconImage) genImageAt(fileName string) error {
	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return err
	}

	opt := truetype.Options{
		Size:              110,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	imageWidth := 128
	imageHeight := 128
	textTopMargin := 105

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	face := truetype.NewFace(ft, &opt)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textIcon.Color),
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	dr.Dot.X = (fixed.I(imageWidth) - dr.MeasureString(*textIcon.Text)) / 2
	dr.Dot.Y = fixed.I(textTopMargin)

	dr.DrawString(*textIcon.Text)

	err = flushRGBA(fileName, img)
	if err != nil {
		return err
	}

	return nil
}

type RectangleIconImage struct {
	Rect  image.Rectangle
	Color *color.RGBA
}

func (rectangleIcon *RectangleIconImage) genImageAt(fileName string) error {
	imageWidth := 100
	imageHeight := 100

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	draw.Draw(
		img,
		rectangleIcon.Rect,
		image.NewUniform(rectangleIcon.Color),
		image.Point{},
		draw.Src,
	)

	err := flushRGBA(fileName, img)
	if err != nil {
		return err
	}

	return nil
}

func flushRGBA(fileName string, rgba *image.RGBA) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	err = png.Encode(buf, rgba)
	if err != nil {
		return err
	}

	_, err = file.Write(buf.Bytes())

	return nil
}
