package nanogui

import (
	"fmt"
	"github.com/goxjs/glfw"
	"github.com/shibukawa/nanovgo"
)

type ButtonFlags int

const (
	NormalButtonType ButtonFlags = 1
	RadioButtonType  ButtonFlags = 2
	ToggleButtonType ButtonFlags = 4
	PopupButtonType  ButtonFlags = 8
)

type ButtonIconPosition int

const (
	Left ButtonIconPosition = iota
	LeftCentered
	Right
	RightCentered
)

type Button struct {
	WidgetImplement

	caption         string
	icon            Icon
	imageIcon       int
	iconPosition    ButtonIconPosition
	pushed          bool
	flags           ButtonFlags
	backgroundColor nanovgo.Color
	textColor       nanovgo.Color
	callback        func()
	changeCallback  func(bool)
	buttonGroup     []*Button
}

func NewButton(parent Widget, labels ...string) *Button {
	var label string
	switch len(labels) {
	case 0:
		label = "Untitled"
	case 1:
		label = labels[0]
	default:
		panic("NewButton can accept only one extra parameter (label)")
	}

	button := &Button{
		caption:      label,
		iconPosition: LeftCentered,
		flags:        NormalButtonType,
	}
	InitWidget(button, parent)
	return button
}

func NewToolButton(parent Widget, icon Icon) *Button {
	button := NewButton(parent, "")
	button.SetCaption("")
	button.SetIcon(icon)
	button.SetFlags(RadioButtonType | ToggleButtonType)
	//button.SetFixedSize(25, 25)
	return button
}

func NewToolButtonByImage(parent Widget, img int) *Button {
	button := NewButton(parent, "")
	button.SetCaption("")
	button.SetImageIcon(img)
	button.SetFlags(RadioButtonType | ToggleButtonType)
	//button.SetFixedSize(25, 25)
	return button
}

func (b *Button) Caption() string {
	return b.caption
}

func (b *Button) SetCaption(caption string) {
	b.caption = caption
}

func (b *Button) BackgroundColor() nanovgo.Color {
	return b.backgroundColor
}

func (b *Button) SetBackgroundColor(c nanovgo.Color) {
	b.backgroundColor = c
}

func (b *Button) TextColor() nanovgo.Color {
	if !b.enabled {
		return b.theme.DisabledTextColor
	} else if b.textColor.A == 0.0 {
		return b.theme.TextColor
	}
	return b.textColor
}

func (b *Button) SetTextColor(c nanovgo.Color) {
	b.textColor = c
}

func (b *Button) Icon() Icon {
	return b.icon
}

func (b *Button) SetIcon(i Icon) {
	b.icon = i
	b.imageIcon = 0
}

func (b *Button) ImageIcon() int {
	return b.imageIcon
}

func (b *Button) SetImageIcon(i int) {
	b.imageIcon = i
	b.icon = 0
}
func (b *Button) Flags() ButtonFlags {
	return b.flags
}

func (b *Button) SetFlags(f ButtonFlags) {
	b.flags = f
}

func (b *Button) IconPosition() ButtonIconPosition {
	return b.iconPosition
}

func (b *Button) SetIconPosition(p ButtonIconPosition) {
	b.iconPosition = p
}

func (b *Button) Pushed() bool {
	return b.pushed
}

func (b *Button) SetPushed(p bool) {
	b.pushed = p
}

// SetCallback set the push callback (for any type of button)
func (b *Button) SetCallback(callback func()) {
	b.callback = callback
}

// SetChangeCallback set the change callback (for toggle buttons)
func (b *Button) SetChangeCallback(callback func(bool)) {
	b.changeCallback = callback
}

// SetButtonGroup set the button group (for radio buttons)
func (b *Button) SetButtonGroup(group []*Button) {
	b.buttonGroup = group
}

// ButtonGroup returns the button group
func (b *Button) ButtonGroup() []*Button {
	return b.buttonGroup
}

func (b *Button) FontSize() int {
	if b.fontSize > 0 {
		return b.fontSize
	}
	return b.theme.ButtonFontSize
}

func (b *Button) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	b.WidgetImplement.MouseButtonEvent(b, x, y, button, down, modifier)

	if button == glfw.MouseButton1 && b.enabled {
		pushedBackup := b.pushed
		if down {
			if b.flags&RadioButtonType != 0 {
				if len(b.buttonGroup) == 0 {
					for _, child := range self.Parent().Children() {
						button, ok := child.(*Button)
						if ok && button != b && button.Flags()&RadioButtonType != 0 && button.Pushed() {
							button.SetPushed(false)
							if button.changeCallback != nil {
								button.changeCallback(false)
							}
						}
					}
				} else {
					for _, button := range b.buttonGroup {
						if button != b && button.Flags()&RadioButtonType != 0 && button.Pushed() {
							button.SetPushed(false)
							if button.changeCallback != nil {
								button.changeCallback(false)
							}
						}
					}
				}
			} else if b.flags&PopupButtonType != 0 {
				for _, widget := range b.Parent().Children() {
					button, ok := widget.(*Button)
					if ok && button != b && button.Flags()&PopupButtonType != 0 && button.Pushed() {
						button.SetPushed(false)
						if button.changeCallback != nil {
							button.changeCallback(false)
						}
					}
				}
			}
			if b.flags&ToggleButtonType != 0 {
				b.pushed = !b.pushed
			} else {
				b.pushed = true
			}
		} else if b.pushed {
			if b.Contains(x, y) && b.callback != nil {
				b.callback()
			}
			if b.flags&NormalButtonType != 0 {
				b.pushed = false
			}
		}
		if pushedBackup != b.pushed && b.changeCallback != nil {
			b.changeCallback(b.pushed)
		}
		return true
	}
	return false
}

func (b *Button) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	fontSize := float32(b.FontSize())

	ctx.SetFontSize(fontSize)
	ctx.SetFontFace(b.theme.FontBold)
	tw, _ := ctx.TextBounds(0, 0, b.caption)
	var iw float32
	ih := fontSize

	if b.icon > 0 {
		ih *= 1.5 / 2
		ctx.SetFontFace(b.theme.FontIcons)
		ctx.SetFontSize(ih)
		iw, _ = ctx.TextBounds(0, 0, string([]rune{rune(b.icon)}))
		iw += float32(b.y) * 0.15
	} else if b.imageIcon > 0 {
		ih *= 0.9
		w, h, _ := ctx.ImageSize(b.imageIcon)
		iw = float32(w) * ih / float32(h)
	}
	return int(tw + iw + 20), int(fontSize) + 10
}

func (b *Button) Draw(ctx *nanovgo.Context) {
	b.WidgetImplement.Draw(ctx)

	bx := float32(b.x)
	by := float32(b.y)
	bw := float32(b.w)
	bh := float32(b.h)

	var gradTop nanovgo.Color
	var gradBot nanovgo.Color

	if b.pushed {
		gradTop = b.theme.ButtonGradientTopPushed
		gradBot = b.theme.ButtonGradientBotPushed
	} else if b.mouseFocus && b.enabled {
		gradTop = b.theme.ButtonGradientTopFocused
		gradBot = b.theme.ButtonGradientBotFocused
	} else {
		gradTop = b.theme.ButtonGradientTopUnfocused
		gradBot = b.theme.ButtonGradientBotUnfocused
	}
	ctx.BeginPath()
	ctx.RoundedRect(bx+1.0, by+1.0, bw-2.0, bh-2.0, float32(b.theme.ButtonCornerRadius-1))

	if b.backgroundColor.A != 0.0 {
		bgColor := b.backgroundColor
		bgColor.A = 1.0
		ctx.SetFillColor(bgColor)
		ctx.Fill()
		if b.pushed {
			gradTop.A = 0.8
			gradBot.A = 0.8
		} else {
			a := 1 - b.backgroundColor.A
			if !b.enabled {
				a = a*0.5 + 0.5
			}
			gradTop.A = a
			gradBot.A = a
		}
	}

	bg := nanovgo.LinearGradient(bx, by, bx, by+bh, gradTop, gradBot)
	ctx.SetFillPaint(bg)
	ctx.Fill()

	ctx.BeginPath()
	var pOff float32 = 0.0
	if b.pushed {
		pOff = 1.0
	}
	ctx.RoundedRect(bx+0.5, by+1.5-pOff, bw-1.0, bh-2+pOff, float32(b.theme.ButtonCornerRadius))
	ctx.SetStrokeColor(b.theme.BorderLight)
	ctx.Stroke()

	ctx.BeginPath()
	ctx.RoundedRect(bx+0.5, by+0.5, bw-1.0, bh-2, float32(b.theme.ButtonCornerRadius))
	ctx.SetStrokeColor(b.theme.BorderDark)
	ctx.Stroke()

	fontSize := float32(b.FontSize())
	ctx.SetFontSize(fontSize)
	ctx.SetFontFace(b.theme.FontBold)
	tw, _ := ctx.TextBounds(0, 0, b.caption)

	centerX := bx + bw*0.5
	centerY := by + bh*0.5
	textPosX := centerX - tw*0.5
	textPosY := centerY - 1.0

	textColor := b.TextColor()
	if b.icon > 0 || b.imageIcon > 0 {
		var iw, ih float32
		if b.icon > 0 {
			ih = fontSize * 1.5 / 2
			ctx.SetFontSize(ih)
			ctx.SetFontFace(b.theme.FontIcons)
			iw, _ = ctx.TextBounds(0, 0, string([]rune{rune(b.icon)}))
		} else if b.imageIcon > 0 {
			ih = fontSize * 0.9
			w, h, _ := ctx.ImageSize(b.imageIcon)
			iw = float32(w) * ih / float32(h)
		}
		if b.caption != "" {
			iw += float32(b.h) * 0.15
		}
		ctx.SetFillColor(textColor)
		ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignMiddle)
		iconPosX := centerX
		iconPosY := centerY - 1

		switch b.iconPosition {
		case LeftCentered:
			iconPosX -= (tw + iw) * 0.5
			textPosX += iw * 0.5
		case RightCentered:
			iconPosX -= iw * 0.5
			textPosX += tw * 0.5
		case Left:
			iconPosX = bx + 8.0
		case Right:
			iconPosX = bx + bw - iw - 8
		}
		if b.icon > 0 {
			ctx.TextRune(iconPosX, iconPosY, []rune{rune(b.icon)})
		} else {
			var eOff float32 = 0.25
			if b.enabled {
				eOff = 0.5
			}
			imgPaint := nanovgo.ImagePattern(iconPosX, iconPosY-ih*0.5, iw, ih, 0, b.imageIcon, eOff)
			ctx.SetFillPaint(imgPaint)
			ctx.Fill()
		}
	}
	ctx.SetFontSize(fontSize)
	ctx.SetFontFace(b.theme.FontBold)
	ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignMiddle)
	ctx.SetFillColor(b.theme.TextColorShadow)
	ctx.Text(textPosX, textPosY, b.caption)
	ctx.SetFillColor(textColor)
	ctx.Text(textPosX, textPosY+1.0, b.caption)
}

func (b *Button) String() string {
	return fmt.Sprintf("Button [%d,%d-%d,%d] - %s", b.x, b.y, b.w, b.h, b.caption)
}
