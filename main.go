package main

import (
	"bytes"
	"fmt"
	"github.com/distatus/battery"
	"github.com/golang/freetype/truetype"
	"github.com/mattn/go-gtk/gtk"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"time"
)

func main() {
	gtk.Init(&os.Args)
	statusIcon := gtk.NewStatusIcon()
	statusIcon.SetTitle("This is a title")
	statusIcon.SetVisible(true)

	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				batteries, err := battery.GetAll()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				if len(batteries) == 0 {
					os.Exit(1)
				}
				bat := batteries[0]
				println(bat.String())

				filename := "/var/tmp/statray_icon.png"
				percentage := int(bat.Current / bat.Full * 100)
				println(percentage)
				err = genIconImage(filename, strconv.Itoa(percentage), color.RGBA{R: 255, G: 255, B: 255, A: 255})
				if err != nil {
					os.Exit(1)
				}

				statusIcon.SetFromFile(filename)
			}
		}
	}()

	gtk.Main()
}

func genIconImage(filename, text string, color color.RGBA) error {
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
		Src:  image.NewUniform(color),
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	dr.Dot.X = (fixed.I(imageWidth) - dr.MeasureString(text)) / 2
	dr.Dot.Y = fixed.I(textTopMargin)

	dr.DrawString(text)

	buf := &bytes.Buffer{}
	err = png.Encode(buf, img)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(buf.Bytes())

	return nil
}
