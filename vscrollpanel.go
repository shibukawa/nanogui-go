package nanogui

import (
	"github.com/shibukawa/glfw"
	"github.com/shibukawa/nanovgo"
	"runtime"
)

type VScrollPanel struct {
	WidgetImplement

	childPreferredHeight int
	scroll               float32
	scrollPosition       float32
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
	child.SetSize(v.w, v.childPreferredHeight)
}

func (v *VScrollPanel) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	if len(v.children) == 0 {
		return 0, 0
	}
	layout := self.Layout()
	if layout != nil {
		w, h := layout.PreferredSize(self, ctx)
		return w + 12, h
	}
	child := v.children[0]
	w, h := child.PreferredSize(child, ctx)
	return w + 12, h
}

func (v *VScrollPanel) MouseDragEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if len(v.children) == 0 {
		return false
	}
	if v.h < v.childPreferredHeight {
		h := float32(v.h)
		ph := float32(v.childPreferredHeight)
		if runtime.GOOS != "darwin" {
			relY = -relY
		}
		scrollAmount := float32(relY) * 2
		v.scrollPosition = clampF(v.scrollPosition-scrollAmount, 0.0, ph-h)
		v.scroll = clampF(v.scrollPosition/(ph-h), 0.0, 1.0)
	} else {
		v.scroll = 0.0
	}
	return true
}

func (v *VScrollPanel) ScrollEvent(self Widget, x, y, relX, relY int) bool {
	if v.h < v.childPreferredHeight {
		h := float32(v.h)
		ph := float32(v.childPreferredHeight)
		scrollAmount := float32(relY) * 2

		v.scrollPosition = clampF(v.scrollPosition-scrollAmount, 0.0, ph-h)
		v.scroll = clampF(v.scrollPosition/(ph-h), 0.0, 1.0)
	} else {
		v.scroll = 0.0
		v.scrollPosition = 0.0
	}
	return true
}

func (v *VScrollPanel) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	if len(v.children) == 0 {
		return false
	}
	child := v.children[0]
	shift := int(v.scroll * float32(v.childPreferredHeight-v.h))
	return child.MouseButtonEvent(child, x, y+shift, button, down, modifier)
}

func (v *VScrollPanel) MouseMotionEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if len(v.children) == 0 {
		return false
	}
	child := v.children[0]
	shift := int(v.scroll * float32(v.childPreferredHeight-v.h))
	return child.MouseMotionEvent(child, x, y+shift, relX, relY, button, modifier)
}

func (v *VScrollPanel) Draw(self Widget, ctx *nanovgo.Context) {
	if len(v.children) == 0 {
		return
	}
	x := float32(v.x)
	y := float32(v.y)
	w := float32(v.w)
	h := float32(v.h)

	child := v.children[0]
	layout := self.Layout()
	if layout != nil {
		_, v.childPreferredHeight = layout.PreferredSize(self, ctx)
	} else {
		_, v.childPreferredHeight = child.PreferredSize(child, ctx)
	}

	ctx.Save()
	ctx.Translate(x, y)
	ctx.Scissor(0, 0, w, h)
	ctx.Translate(0, -v.scroll*(float32(v.childPreferredHeight)-h))
	if child.Visible() {
		child.Draw(child, ctx)
	}
	ctx.Restore()
	if v.childPreferredHeight > v.h {
		scrollH := h * minF(1.0, h/float32(v.childPreferredHeight))
		scrollH = minF(maxF(20.0, scrollH), h)
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
}

func (v *VScrollPanel) IsClipped(x, y, w, h int) bool {
	scroll := int(v.scroll * (float32(v.childPreferredHeight) - float32(v.h)))
	return v.Parent().IsClipped(x+v.x, y-scroll+v.y, w, h)
}

func (v *VScrollPanel) String() string {
	return v.StringHelper("VScrollPanel", "")
}

type VScrollPanelChild struct {
	WidgetImplement
}

func NewVScrollPanelChild(parent *VScrollPanel) *VScrollPanelChild {
	result := &VScrollPanelChild{}
	InitWidget(result, parent)
	return result
}

func (v *VScrollPanelChild) Size() (int, int) {
	return v.Parent().Size()
}

func (v *VScrollPanelChild) SetSize(w, h int) {
	if h > 400 {
		v.Parent().Parent().SetSize(w, 400)
		v.Parent().SetSize(w, 400)
	} else {
		v.Parent().Parent().SetSize(w, h)
		v.Parent().SetSize(w, h)
	}
}

func (v *VScrollPanelChild) Width() int {
	return v.Parent().Width()
}

func (v *VScrollPanelChild) SetWidth(w int) {
	v.Parent().Parent().SetWidth(w)
	v.Parent().SetWidth(w)
}

func (v *VScrollPanelChild) Height() int {
	return v.Parent().Height()
}

func (v *VScrollPanelChild) SetHeight(h int) {
	v.Parent().Parent().SetHeight(h)
	v.Parent().SetHeight(h)
}

func (v *VScrollPanelChild) FixedSize() (int, int) {
	return v.Parent().Parent().FixedSize()
}

func (v *VScrollPanelChild) SetFixedSize(w, h int) {
	v.Parent().Parent().SetFixedSize(w, h)
}

func (v *VScrollPanelChild) FixedWidth() int {
	return v.Parent().Parent().FixedWidth()
}

func (v *VScrollPanelChild) FixedHeight() int {
	return v.Parent().Parent().FixedHeight()
}

func (v *VScrollPanelChild) SetFixedWidth(w int) {
	v.Parent().Parent().SetFixedWidth(w)
}

func (v *VScrollPanelChild) SetFixedHeight(h int) {
	v.Parent().Parent().SetFixedHeight(h)
}

func (v *VScrollPanelChild) Visible() bool {
	return v.Parent().Parent().Visible()
}

func (v *VScrollPanelChild) SetVisible(flag bool) {
	v.Parent().Parent().SetVisible(flag)
}

func (v *VScrollPanelChild) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	v.WidgetImplement.MouseButtonEvent(self, x, y, button, down, modifier)
	return true
}

func (v *VScrollPanelChild) MouseMotionEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	v.WidgetImplement.MouseMotionEvent(self, x, y, relX, relY, button, modifier)
	return true
}

func (v *VScrollPanelChild) ScrollEvent(self Widget, x, y, relX, relY int) bool {
	v.WidgetImplement.ScrollEvent(self, x, y, relX, relY)
	return true
}

func (v *VScrollPanelChild) String() string {
	return v.StringHelper("VScrollPanelChild", "")
}

func (v *VScrollPanelChild) IsClipped(cx, cy, cw, ch int) bool {
	parent := v.Parent()
	if parent == nil {
		return false
	}
	return parent.IsClipped(cx+v.x, cy+v.y, cw, ch)
}
