package icon

import (
	"fmt"
	"github.com/mattn/go-gtk/gtk"
	"github.com/shirou/gopsutil/cpu"
	"image"
	"image/color"
	"time"
)

const cpuIconPath = "/var/tmp/statray_icon_cpu_%d.png"

func NewCpuIcon() (AnimateIcon, error) {
	coreCount, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	gtkIcons := make([]*gtk.StatusIcon, coreCount)
	for core := range gtkIcons {
		gtkIcons[core] = gtk.NewStatusIcon()
	}

	return &cpuIcon{
		GtkIcons:  gtkIcons,
		CoreCount: coreCount,
	}, nil
}

type cpuIcon struct {
	GtkIcons  []*gtk.StatusIcon
	CoreCount int
}

func (icon *cpuIcon) Update() error {
	percents, err := cpu.Percent(time.Duration(0), true)
	if err != nil {
		return nil
	}

	for iconIdx, gtkIcon := range icon.GtkIcons {
		fileName := fmt.Sprintf(cpuIconPath, iconIdx)
		percentage := percents[iconIdx]

		var clr color.RGBA
		switch {
		case percentage >= 90:
			clr = colorRed
		default:
			clr = colorWhite
		}

		err := genRectangleIconImage(
			RectangleIconImage{
				Rect: image.Rectangle{
					Min: image.Point{X: 15, Y: 85 - (int(percentage * 70 / 100))},
					Max: image.Point{X: 85, Y: 85},
				},
				Color: &clr,
			},
			fileName,
		)
		if err != nil {
			return err
		}

		gtkIcon.SetFromFile(fileName)
		gtkIcon.SetVisible(true)
	}

	return nil
}
