package nanogui

import (
	"fmt"
	"github.com/shibukawa/nanovgo"
)

type ProgressBar struct {
	WidgetImplement

	value float32
}

func NewProgressBar(parent Widget) *ProgressBar {
	progressBar := &ProgressBar{}
	InitWidget(progressBar, parent)
	return progressBar
}

func (p *ProgressBar) Value() float32 {
	return p.value
}

func (p *ProgressBar) SetValue(value float32) {
	p.value = value
}

func (p *ProgressBar) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	return 70, 12
}

func (p *ProgressBar) Draw(ctx *nanovgo.Context) {
	px := float32(p.x)
	py := float32(p.y)
	pw := float32(p.w)
	ph := float32(p.h)
	p.WidgetImplement.Draw(ctx)
	paint := nanovgo.BoxGradient(px+1, py+1, pw-2, ph, 3, 4, nanovgo.MONO(0, 32), nanovgo.MONO(0, 92))
	ctx.BeginPath()
	ctx.RoundedRect(px, py, pw, ph, 3)
	ctx.SetFillPaint(paint)
	ctx.Fill()

	value := clampF(p.value, 0.0, 1.0)
	barPos := (pw - 2) * value
	barPaint := nanovgo.BoxGradient(px, py, barPos+1.5, ph-1, 3, 4, nanovgo.MONO(220, 100), nanovgo.MONO(128, 100))
	ctx.BeginPath()
	ctx.RoundedRect(px+1, py+1, barPos, ph-2, 3)
	ctx.SetFillPaint(barPaint)
	ctx.Fill()
}

func (p *ProgressBar) String() string {
	return fmt.Sprintf("ProgressBar [%d,%d-%d,%d] - %f", p.x, p.y, p.w, p.h, p.value)
}
