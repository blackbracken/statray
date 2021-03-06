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

	var iconColor color.RGBA
	switch {
	case onFullCharge:
		fallthrough
	case bat.State == battery.Charging:
		iconColor = colorYellow
	case percentage >= 80:
		iconColor = colorGreen
	case percentage <= 20:
		iconColor = colorRed
	default:
		iconColor = colorWhite
	}

	textIconImg :=
		TextIconImage{
			Text:  &iconText,
			Color: &iconColor,
		}
	err = textIconImg.genImageAt(batteryIconPath)
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
