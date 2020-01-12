package icon

import (
	"errors"
	"fmt"
	"github.com/distatus/battery"
	"github.com/mattn/go-gtk/gtk"
	"image/color"
	"math"
	"strconv"
)

const batteryIconPath = "/var/tmp/statray_icon_battery.png"

func NewBatteryIcon() AnimateIcon {
	statusIcon := gtk.NewStatusIcon()
	statusIcon.SetTitle("An icon to show status of your battery")

	return &batteryIcon{GtkIcon: statusIcon}
}

type batteryIcon struct {
	GtkIcon *gtk.StatusIcon
}

func (icon *batteryIcon) Update() error {
	bat, err := getSingleBattery()
	if err != nil {
		return err
	}

	percentage := int(math.Round(bat.Current / bat.Full * 100))
	onFullCharge := percentage >= 100

	var iconText string
	if onFullCharge {
		iconText = "F"
	} else {
		iconText = strconv.Itoa(percentage)
	}

	println(bat.String())

	var iconColor color.RGBA
	switch {
	case onFullCharge:
		fallthrough
	case bat.State == battery.Charging:
		iconColor = color.RGBA{R: 242, G: 211, B: 36, A: 255}
	case percentage >= 80:
		iconColor = color.RGBA{R: 73, G: 204, B: 130, A: 255}
	case percentage <= 20:
		iconColor = color.RGBA{R: 227, G: 78, B: 73, A: 255}
	default:
		iconColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}

	err = genTextIconImage(
		TextIconImage{
			Text:  &iconText,
			Color: &iconColor,
		},
		batteryIconPath,
	)
	if err != nil {
		return err
	}
	icon.GtkIcon.SetFromFile(batteryIconPath)

	icon.GtkIcon.SetTooltipText(
		fmt.Sprintf("Capacity: %d%% ( %.1f / %.1f [Wh] )", percentage, bat.Current/1000.0, bat.Full/1000.0))
	icon.GtkIcon.SetVisible(true)

	return nil
}

func getSingleBattery() (*battery.Battery, error) {
	batteries, err := battery.GetAll()
	if err != nil || len(batteries) == 0 {
		return nil, errors.New("Failed to find a battery. ")
	}

	return batteries[0], nil
}
