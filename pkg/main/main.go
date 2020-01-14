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

	cpuIcon, err := icon.NewCpuIcon()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	animateIcons := []icon.AnimateIcon{
		icon.NewBatteryIcon(),
		cpuIcon,
	}
	updateIcons(animateIcons)

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
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
