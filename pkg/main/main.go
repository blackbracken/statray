package main

import (
	"fmt"
	"github.com/blackbracken/statray/pkg/icon"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"time"
)

func main() {
	gtk.Init(&os.Args)

	animateIcons := []icon.AnimateIcon{
		icon.NewBatteryIcon(),
	}
	updateIcons(animateIcons)

	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				updateIcons(animateIcons)
			}
		}
	}()

	gtk.Main()
}

func updateIcons(animateIcons []icon.AnimateIcon) {
	for _, animateIcon := range animateIcons {
		err := animateIcon.Update()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
