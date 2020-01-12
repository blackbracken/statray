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

type AnimateIcon interface {
	Update() error
}

type TextIconImage struct {
	Text  *string
	Color *color.RGBA
}

type RectangleIconImage struct {
	Rect  image.Rectangle
	Color *color.RGBA
}

func genTextIconImage(textIcon TextIconImage, fileName string) error {
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

func genRectangleIconImage(rectangleIcon RectangleIconImage, fileName string) error {
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
