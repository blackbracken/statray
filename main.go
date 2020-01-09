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
	"image/png"
	"os"
	"time"
)

func main() {
	fmt.Println("Hello, World!")

	gtk.Init(&os.Args)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Hello Go-GTK")
	window.Connect("destroy", gtk.MainQuit)

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case tick := <-ticker.C:
				fmt.Println("Tick at", tick)
				batteries, err := battery.GetAll()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				for i, bat := range batteries {
					fmt.Println(i)
					window.SetTitle(fmt.Sprint(bat.State))
				}
			}
		}
	}()

	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	opt := truetype.Options{
		Size:              90,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	imageWidth := 100
	imageHeight := 100
	textTopMargin := 90
	text := "A"

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	face := truetype.NewFace(ft, &opt)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	dr.Dot.X = (fixed.I(imageWidth) - dr.MeasureString(text)) / 2
	dr.Dot.Y = fixed.I(textTopMargin)

	dr.DrawString(text)

	buf := &bytes.Buffer{}
	err = png.Encode(buf, img)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	file, err := os.Create(`test.png`)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	file.Write(buf.Bytes())

	statusIcon := gtk.NewStatusIconFromFile(`test.png`)
	statusIcon.SetTitle("This is a title")
	statusIcon.SetVisible(true)

	window.ShowAll()
	window.SetSizeRequest(400, 250)

	gtk.Main()
}
