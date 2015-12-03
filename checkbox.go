package nanogui

import (
	"fmt"
	"github.com/goxjs/glfw"
	"github.com/shibukawa/nanovgo"
)

type CheckBox struct {
	WidgetImplement
	caption  string
	pushed   bool
	checked  bool
	callback func(bool)
}

func NewCheckBox(parent Widget, caption string) *CheckBox {
	if caption == "" {
		caption = "Untitled"
	}
	checkBox := &CheckBox{
		caption: caption,
	}
	InitWidget(checkBox, parent)
	return checkBox
}

func (c *CheckBox) Caption() string {
	return c.caption
}

func (c *CheckBox) SetCaption(caption string) {
	c.caption = caption
}

func (c *CheckBox) Checked() bool {
	return c.checked
}

func (c *CheckBox) SetChecked(checked bool) {
	c.checked = checked
}

func (c *CheckBox) Pushed() bool {
	return c.pushed
}

func (c *CheckBox) SetPushed(pushed bool) {
	c.pushed = pushed
}

func (c *CheckBox) SetCallback(callback func(bool)) {
	c.callback = callback
}

func (c *CheckBox) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	c.WidgetImplement.MouseButtonEvent(self, x, y, button, down, modifier)
	if !c.enabled {
		return false
	}
	if button == glfw.MouseButton1 {
		if down {
			c.pushed = true
		} else if c.pushed {
			if c.Contains(x, y) {
				c.checked = !c.checked
				if c.callback != nil {
					c.callback(c.checked)
				}
			}
			c.pushed = false
		}
		return true
	}
	return false
}

func (c *CheckBox) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	fw, fh := c.FixedSize()
	if fw > 0 || fh > 0 {
		return fw, fh
	}
	fontSize := float32(c.FontSize())
	ctx.SetFontSize(fontSize)
	ctx.SetFontFace(c.theme.FontNormal)
	w, _ := ctx.TextBounds(0, 0, c.caption)
	return int(w + 1.7*fontSize), int(fontSize * 1.3)
}

func (c *CheckBox) Draw(ctx *nanovgo.Context) {
	cx := float32(c.x)
	cy := float32(c.y)
	ch := float32(c.h)
	c.WidgetImplement.Draw(ctx)
	fontSize := float32(c.FontSize())
	ctx.SetFontSize(fontSize)
	ctx.SetFontFace(c.theme.FontNormal)
	if c.enabled {
		ctx.SetFillColor(c.theme.TextColor)
	} else {
		ctx.SetFillColor(c.theme.DisabledTextColor)
	}
	ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignMiddle)
	ctx.Text(cx+1.2*ch+5, cy+ch*0.5, c.caption)
	var bgAlpha uint8
	if c.pushed {
		bgAlpha = 100
	} else {
		bgAlpha = 32
	}
	bgPaint := nanovgo.BoxGradient(cx+1.5, cy+1.5, ch-2.0, ch-2.0, 3, 3, nanovgo.MONO(0, bgAlpha), nanovgo.MONO(0, 180))
	ctx.BeginPath()
	ctx.RoundedRect(cx+1.0, cy+1.0, ch-2.0, ch-2.0, 3)
	ctx.SetFillPaint(bgPaint)
	ctx.Fill()

	if c.checked {
		ctx.SetFontSize(ch)
		ctx.SetFontFace(c.theme.FontIcons)
		if c.enabled {
			ctx.SetFillColor(c.theme.IconColor)
		} else {
			ctx.SetFillColor(c.theme.DisabledTextColor)
		}
		ctx.SetTextAlign(nanovgo.AlignCenter | nanovgo.AlignMiddle)
		ctx.Text(cx+ch*0.5+1.0, cy+ch*0.5, string([]rune{rune(IconCheck)}))
	}
}

func (c *CheckBox) String() string {
	return fmt.Sprintf("CheckBox [%d,%d-%d,%d] - %s", c.x, c.y, c.w, c.h, c.caption)
}
