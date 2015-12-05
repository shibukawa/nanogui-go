package nanogui

import (
	"fmt"
	"github.com/shibukawa/nanovgo"
)

type ColorPicker struct {
	PopupButton

	callback   func(color nanovgo.Color)
	colorWheel *ColorWheel
	pickButton *Button
}

func NewColorPicker(parent Widget, colors ...nanovgo.Color) *ColorPicker {
	var color nanovgo.Color
	switch len(colors) {
	case 0:
		color = nanovgo.RGBAf(1.0, 0.0, 0.0, 1.0)
	case 1:
		color = colors[0]
	default:
		panic("NewColorPicker can accept only one extra parameter (color)")
	}

	colorPicker := &ColorPicker{}

	// init PopupButton member
	colorPicker.chevronIcon = IconRightOpen
	colorPicker.SetIconPosition(ButtonIconLeftCentered)
	colorPicker.SetFlags(ToggleButtonType | PopupButtonType)
	parentWindow := parent.FindWindow()

	colorPicker.popup = NewPopup(parentWindow.Parent(), parentWindow)
	colorPicker.popup.SetLayout(NewGroupLayout())

	colorPicker.colorWheel = NewColorWheel(colorPicker.popup)

	colorPicker.pickButton = NewButton(colorPicker.popup, "Pick")
	colorPicker.pickButton.SetFixedSize(100, 25)

	InitWidget(colorPicker, parent)

	colorPicker.SetColor(color)

	colorPicker.PopupButton.SetChangeCallback(func(flag bool) {
		colorPicker.SetColor(colorPicker.BackgroundColor())
		if colorPicker.callback != nil {
			colorPicker.callback(colorPicker.BackgroundColor())
		}
	})

	colorPicker.colorWheel.SetCallback(func(color nanovgo.Color) {
		colorPicker.pickButton.SetBackgroundColor(color)
		colorPicker.pickButton.SetTextColor(color.ContrastingColor())
	})

	colorPicker.pickButton.SetCallback(func() {
		color := colorPicker.colorWheel.Color()
		colorPicker.SetPushed(false)
		colorPicker.SetColor(color)
		if colorPicker.callback != nil {
			colorPicker.callback(colorPicker.BackgroundColor())
		}
	})

	return colorPicker
}

func (c *ColorPicker) SetCallback(callback func(color nanovgo.Color)) {
	c.callback = callback
}

func (c *ColorPicker) Color() nanovgo.Color {
	return c.BackgroundColor()
}

func (c *ColorPicker) SetColor(color nanovgo.Color) {
	if !c.pushed {
		fgColor := color.ContrastingColor()
		c.SetBackgroundColor(color)
		c.SetTextColor(fgColor)
		c.colorWheel.SetColor(color)
		c.pickButton.SetBackgroundColor(color)
		c.pickButton.SetTextColor(fgColor)
	}
}

func (c *ColorPicker) String() string {
	cw := c.colorWheel
	return fmt.Sprintf("ColorPicker [%d,%d-%d,%d] - h:%f s:%f l:%f", c.x, c.y, c.w, c.h, cw.hue, cw.saturation, cw.lightness)
}
