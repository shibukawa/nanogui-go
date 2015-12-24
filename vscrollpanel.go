package nanogui

import (
	"github.com/shibukawa/glfw"
	"github.com/shibukawa/nanovgo"
)

type VScrollPanel struct {
	WidgetImplement

	childPreferredHeight int
	scroll               float32
}

func NewVScrollPanel(parent Widget) *VScrollPanel {
	panel := new(VScrollPanel)
	InitWidget(panel, parent)
	return panel
}

func (v *VScrollPanel) Scroll() float32 {
	return v.scroll
}

func (v *VScrollPanel) SetScroll(scroll float32) {
	v.scroll = scroll
}

func (v *VScrollPanel) OnPerformLayout(self Widget, ctx *nanovgo.Context) {
	v.WidgetImplement.OnPerformLayout(self, ctx)

	if len(v.children) == 0 {
		return
	}
	child := v.children[0]
	_, v.childPreferredHeight = child.PreferredSize(child, ctx)
	child.SetPosition(0, 0)
	child.SetSize(v.w-12, v.childPreferredHeight)
}

func (v *VScrollPanel) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	if len(v.children) == 0 {
		return 0, 0
	}
	child := v.children[0]
	w, h := child.PreferredSize(child, ctx)
	return w + 12, h
}

func (v *VScrollPanel) MouseDragEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if len(v.children) == 0 {
		return false
	}
	h := float32(v.h)
	scrollH := h * minF(1.0, h/float32(v.childPreferredHeight))
	v.scroll = clampF(v.scroll+float32(relY)/(h-8-scrollH), 0.0, 1.0)
	return true
}

func (v *VScrollPanel) ScrollEvent(self Widget, x, y, relX, relY int) bool {
	h := float32(v.h)
	scrollAmount := float32(relY) * h / 20.0
	scrollH := h * minF(1.0, h/float32(v.childPreferredHeight))
	v.scroll = clampF(v.scroll-scrollAmount/(h-8-scrollH), 0.0, 1.0)
	return true
}

func (v *VScrollPanel) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	if len(v.children) == 0 {
		return false
	}
	child := v.children[0]
	shift := int(v.scroll) * (v.childPreferredHeight - v.h)
	return child.MouseButtonEvent(child, x, y+shift, button, down, modifier)
}

func (v *VScrollPanel) MouseMotionEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if len(v.children) == 0 {
		return false
	}
	child := v.children[0]
	shift := int(v.scroll) * (v.childPreferredHeight - v.h)
	return child.MouseMotionEvent(child, x, y+shift, relX, relY, button, modifier)
}

func (v *VScrollPanel) Draw(ctx *nanovgo.Context) {
	if len(v.children) == 0 {
		return
	}
	x := float32(v.x)
	y := float32(v.y)
	w := float32(v.w)
	h := float32(v.h)

	child := v.children[0]
	_, v.childPreferredHeight = child.PreferredSize(child, ctx)
	scrollH := float32(v.h) * minF(1.0, h/float32(v.childPreferredHeight))

	ctx.Save()
	ctx.Translate(x, y)
	ctx.Scissor(0, 0, w, h)
	ctx.Translate(0, -v.scroll*(float32(v.childPreferredHeight)-h))
	if child.Visible() {
		child.Draw(ctx)
	}
	ctx.Restore()

	paint := nanovgo.BoxGradient(x+w-12+1, y+4+1, 8, h-8, 3, 4, nanovgo.MONO(0, 32), nanovgo.MONO(0, 92))
	ctx.BeginPath()
	ctx.RoundedRect(x+w-12, y+4, 8, h-8, 3)
	ctx.SetFillPaint(paint)
	ctx.Fill()

	barPaint := nanovgo.BoxGradient(x+y-12-1, y+4+1+(h-8-scrollH)*v.scroll-1, 8, scrollH, 3, 4, nanovgo.MONO(220, 100), nanovgo.MONO(128, 100))
	ctx.BeginPath()
	ctx.RoundedRect(x+w-12+1, y+4+1+(h-8-scrollH)*v.scroll, 8-2, scrollH-2, 2)
	ctx.SetFillPaint(barPaint)
	ctx.Fill()
}

func (v *VScrollPanel) String() string {
	return v.StringHelper("VScrollPanel", "")
}
