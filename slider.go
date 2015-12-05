package nanogui

import (
	"fmt"
	"github.com/goxjs/glfw"
	"github.com/shibukawa/nanovgo"
)

type Slider struct {
	WidgetImplement

	value            float32
	highlightColor   nanovgo.Color
	highlightedRange [2]float32
	callback         func(float32)
	finalCallback    func(float32)
}

func NewSlider(parent Widget) *Slider {
	slider := &Slider{}
	InitWidget(slider, parent)
	return slider
}

func (s *Slider) Value() float32 {
	return s.value
}

func (s *Slider) SetValue(v float32) {
	s.value = v
}

func (s *Slider) HighlightColor() nanovgo.Color {
	return s.highlightColor
}

func (s *Slider) SetHighlightColor(c nanovgo.Color) {
	s.highlightColor = c
}

func (s *Slider) HighlightedRange() (float32, float32) {
	return s.highlightedRange[0], s.highlightedRange[1]
}

func (s *Slider) SetHighlightedRange(l, h float32) {
	s.highlightedRange[0] = l
	s.highlightedRange[1] = h
}

func (s *Slider) SetCallback(callback func(float32)) {
	s.callback = callback
}

func (s *Slider) SetFinalCallback(callback func(float32)) {
	s.finalCallback = callback
}

func (s *Slider) MouseDragEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if !s.enabled {
		return false
	}
	s.value = clampF(float32(x-s.x)/float32(s.w), 0.0, 1.0)
	if s.callback != nil {
		s.callback(s.value)
	}
	return true
}

func (s *Slider) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	if !s.enabled {
		return false
	}
	s.value = clampF(float32(x-s.x)/float32(s.w), 0.0, 1.0)
	if s.callback != nil {
		s.callback(s.value)
	}
	if s.finalCallback != nil {
		s.finalCallback(s.value)
	}
	return true
}

func (s *Slider) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	return 70, 12
}

func (s *Slider) Draw(ctx *nanovgo.Context) {
	sx := float32(s.x)
	sy := float32(s.y)
	sw := float32(s.w)
	sh := float32(s.h)
	cy := sy + sh*0.5
	kx := sx + s.value*sw
	ky := cy + 0.5
	kr := sh * 0.5

	var a1, a2, a3 uint8
	if s.enabled {
		a1 = 32
		a2 = 128
		a3 = 255
	} else {
		a1 = 10
		a2 = 210
		a3 = 100
	}
	background := nanovgo.BoxGradient(sx, cy-3+1, sw, 6, 3, 3, nanovgo.MONO(0, a1), nanovgo.MONO(0, a2))

	ctx.BeginPath()
	ctx.RoundedRect(sx, cy-3+1, sw, 6, 2)
	ctx.SetFillPaint(background)
	ctx.Fill()

	if s.highlightedRange[0] != s.highlightedRange[1] {
		ctx.BeginPath()
		ctx.RoundedRect(sx+s.highlightedRange[0]*sw, cy-3+1, sw*(s.highlightedRange[1]-s.highlightedRange[0]), 6, 2)
		ctx.SetFillColor(s.highlightColor)
		ctx.Fill()
	}

	knobShadow := nanovgo.RadialGradient(kx, ky, kr-3, kr+3, nanovgo.MONO(0, 64), s.theme.Transparent)
	ctx.BeginPath()
	ctx.Rect(kx-kr-5, ky-kr-5, kr*2+10, kr*2+10+3)
	ctx.Circle(kx, ky, kr)
	ctx.PathWinding(nanovgo.Hole)
	ctx.SetFillPaint(knobShadow)
	ctx.Fill()

	knobPaint := nanovgo.LinearGradient(sx, cy-kr, sx, cy+kr, s.theme.BorderLight, s.theme.BorderMedium)
	knobReversePaint := nanovgo.LinearGradient(sx, cy-kr, sx, cy+kr, s.theme.BorderMedium, s.theme.BorderLight)

	ctx.BeginPath()
	ctx.Circle(kx, ky, kr)
	ctx.SetStrokeColor(s.theme.BorderDark)
	ctx.SetFillPaint(knobPaint)
	ctx.Stroke()
	ctx.Fill()

	ctx.BeginPath()
	ctx.Circle(kx, ky, kr/2)
	ctx.SetStrokePaint(knobReversePaint)
	ctx.SetFillColor(nanovgo.MONO(150, a3))
	ctx.Stroke()
	ctx.Fill()
}

func (s *Slider) String() string {
	return fmt.Sprintf("Slider [%d,%d-%d,%d] - %f", s.x, s.y, s.w, s.h, s.value)
}
