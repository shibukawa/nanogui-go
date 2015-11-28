package nanogui

import (
	"fmt"
	"github.com/goxjs/glfw"
	"github.com/shibukawa/nanovgo"
)

type Window struct {
	WidgetImplement
	title       string
	buttonPanel Widget
	modal       bool
	drag        bool
}

type IWindow interface {
	Widget
	RefreshRelativePlacement()
}

func NewWindow(parent Widget, title string) *Window {
	if title == "" {
		title = "Untitled"
	}
	window := &Window{
		title: title,
	}
	InitWidget(window, parent)
	return window
}

// Title() returns the window title
func (w *Window) Title() string {
	return w.title
}

// SetTitle() sets the window title
func (w *Window) SetTitle(title string) {
	w.title = title
}

// Modal() returns is this a model dialog?
func (w *Window) Modal() bool {
	return w.modal
}

// SetModal() set whether or not this is a modal dialog
func (w *Window) SetModal(m bool) {
	w.modal = m
}

func (w *Window) ButtonPanel() Widget {
	if w.buttonPanel == nil {
		widget := &WidgetImplement{}
		InitWidget(w, widget)
		w.buttonPanel = widget
		w.buttonPanel.SetLayout(NewBoxLayout(Horizontal, Middle, 0, 4))
	}
	return w.buttonPanel
}

// Dispose() disposes the window
func (w *Window) Dispose() {
	var widget Widget = w
	var parent Widget = widget.Parent()
	for parent != nil {
		widget = parent
		parent = widget.Parent()
	}
	screen := parent.(*Screen)
	screen.DisposeWindow(w)
}

// Center() makes the window center in the current Screen
func (w *Window) Center() {
	var widget Widget = w
	var parent Widget = widget.Parent()
	for parent != nil {
		widget = parent
		parent = widget.Parent()
	}
	screen := parent.(*Screen)
	screen.CenterWindow(w)
}

// RefreshRelativePlacement is internal helper function to maintain nested window position values; overridden in \ref Popup
func (w *Window) RefreshRelativePlacement() {
	// overridden in Popup
}

func (w *Window) MouseButtonEvent(x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	if w.WidgetImplement.MouseButtonEvent(x, y, button, down, modifier) {
		return true
	}
	if button == glfw.MouseButton1 {
		w.drag = down && (y-w.y) < w.theme.WindowHeaderHeight
		return true
	}
	return false
}

func (w *Window) MouseDragEvent(x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if w.drag && (button&1<<uint(glfw.MouseButton1)) != 0 {
		pW, pH := w.Parent().Size()
		w.x = clampI(w.x+relX, 0, pW-w.w)
		w.y = clampI(w.y+relY, 0, pH-w.h)
		return true
	}
	return false
}

func (w *Window) ScrollEvent(x, y, relX, relY int) bool {
	w.WidgetImplement.ScrollEvent(x, y, relX, relY)
	return true
}

func (w *Window) PreferredSize(ctx *nanovgo.Context, widget Widget) (int, int) {
	if w.buttonPanel != nil {
		w.buttonPanel.SetVisible(false)
	}
	width, height := w.WidgetImplement.PreferredSize(ctx, widget)
	if w.buttonPanel != nil {
		w.buttonPanel.SetVisible(true)
	}
	ctx.SetFontSize(18.0)
	ctx.SetFontFace(w.theme.FontBold)
	_, bounds := ctx.TextBounds(0, 0, w.title)

	return maxI(width, int(bounds[2]-bounds[0])+20), maxI(height, int(bounds[3]-bounds[1]))
}

func (w *Window) OnPerformLayout(ctx *nanovgo.Context, widget Widget) {
	if w.buttonPanel == nil {
		w.WidgetImplement.OnPerformLayout(ctx, widget)
	} else {
		w.buttonPanel.SetVisible(false)
		w.WidgetImplement.OnPerformLayout(ctx, w)
		for _, c := range w.buttonPanel.Children() {
			c.SetFixedSize(22, 22)
			c.SetFontSize(15)
		}
		w.buttonPanel.SetVisible(true)
		w.buttonPanel.SetSize(w.Width(), 22)
		panelW, _ := w.buttonPanel.PreferredSize(ctx, w.buttonPanel)
		w.buttonPanel.SetPosition(w.Width()-(panelW+5), 3)
		w.buttonPanel.OnPerformLayout(ctx, w.buttonPanel)
	}
}

func (w *Window) Draw(ctx *nanovgo.Context) {
	ds := float32(w.theme.WindowDropShadowSize)
	cr := float32(w.theme.WindowCornerRadius)
	hh := float32(w.theme.WindowHeaderHeight)

	// Draw window
	wx := float32(w.x)
	wy := float32(w.y)
	ww := float32(w.w)
	wh := float32(w.h)
	ctx.Save()
	ctx.BeginPath()
	ctx.RoundedRect(wx, wy, ww, wh, cr)
	if w.mouseFocus {
		ctx.SetFillColor(w.theme.WindowFillFocused)
	} else {
		ctx.SetFillColor(w.theme.WindowFillUnfocused)
	}
	ctx.Fill()

	// Draw a drop shadow
	shadowPaint := nanovgo.BoxGradient(wx, wy, ww, wh, cr*2, ds*2, w.theme.DropShadow, w.theme.Transparent)
	ctx.BeginPath()
	ctx.Rect(wx-ds, wy-ds, ww+ds*2, wh+ds*2)
	ctx.RoundedRect(wx, wy, ww, wh, cr)
	ctx.PathWinding(nanovgo.Hole)
	ctx.SetFillPaint(shadowPaint)
	ctx.Fill()

	if w.title != "" {
		headerPaint := nanovgo.LinearGradient(wx, wy, ww, wh+hh, w.theme.WindowHeaderGradientTop, w.theme.WindowHeaderGradientBot)

		ctx.BeginPath()
		ctx.RoundedRect(wx, wy, ww, hh, cr)
		ctx.SetFillPaint(headerPaint)
		ctx.Fill()

		ctx.BeginPath()
		ctx.RoundedRect(wx, wy, ww, wh, cr)
		ctx.SetStrokeColor(w.theme.WindowHeaderSepTop)
		ctx.Scissor(wx, wy, ww, 0.5)
		ctx.Stroke()
		ctx.ResetScissor()

		ctx.BeginPath()
		ctx.MoveTo(wx+0.5, wy+hh-1.5)
		ctx.LineTo(wx+ww-0.5, wy+hh-1.5)
		ctx.SetStrokeColor(w.theme.WindowHeaderSepTop)
		ctx.Stroke()

		ctx.SetFontSize(18.0)
		ctx.SetFontFace(w.theme.FontBold)
		ctx.SetTextAlign(nanovgo.AlignCenter | nanovgo.AlignMiddle)
		ctx.SetFontBlur(2.0)
		ctx.SetFillColor(w.theme.DropShadow)
		ctx.Text(wx+ww*0.5, wy+hh*0.5, w.title)
		ctx.SetFontBlur(0.0)
		if w.focused {
			ctx.SetFillColor(w.theme.WindowTitleFocused)
		} else {
			ctx.SetFillColor(w.theme.WindowTitleUnfocused)
		}
		ctx.Text(wx+ww*0.5, wy+hh*0.5-1, w.title)
	}
	ctx.Restore()
	w.WidgetImplement.Draw(ctx)
}

func (w *Window) FindWindow() *Window {
	return w
}

func (w *Window) String() string {
	return fmt.Sprintf("Window [%d,%d-%d,%d] - %s", w.x, w.y, w.w, w.h, w.title)
}
