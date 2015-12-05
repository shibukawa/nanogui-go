package nanogui

import (
	"fmt"
	"github.com/goxjs/glfw"
	"github.com/shibukawa/nanovgo"
)

// Widget is base class of all widgets
//
// Widget is the base class of all widgets in nanogui. It can
// also be used as an panel to arrange an arbitrary number of child
// widgets using a layout generator (see Layout)
type Widget interface {
	Parent() Widget
	SetParent(parent Widget)
	Layout() Layout
	SetLayout(layout Layout)
	Theme() *Theme
	SetTheme(theme *Theme)
	Position() (int, int)
	SetPosition(x, y int)
	AbsolutePosition() (int, int)
	Size() (int, int)
	SetSize(w, h int)
	Width() int
	SetWidth(w int)
	Height() int
	SetHeight(h int)
	FixedSize() (int, int)
	SetFixedSize(w, h int)
	FixedWidth() int
	SetFixedWidth(w int)
	FixedHeight() int
	SetFixedHeight(h int)
	Visible() bool
	SetVisible(v bool)
	VisibleRecursive() bool
	ChildCount() int
	Children() []Widget
	AddChild(self, w Widget)
	RemoveChildByIndex(i int)
	RemoveChild(w Widget)
	FindWindow() IWindow
	SetID(id string)
	ID() string
	Enabled() bool
	SetEnabled(e bool)
	Focused() bool
	SetFocused(f bool)
	RequestFocus(self Widget)
	Tooltip() string
	SetTooltip(s string)
	FontSize() int
	SetFontSize(s int)
	HasFontSize() bool
	Cursor() Cursor
	SetCursor(c Cursor)
	Contains(x, y int) bool
	FindWidget(self Widget, x, y int) Widget
	MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool
	MouseMotionEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool
	MouseDragEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool
	MouseEnterEvent(self Widget, x, y int, enter bool) bool
	ScrollEvent(self Widget, x, y, relX, relY int) bool
	FocusEvent(self Widget, f bool) bool
	KeyboardEvent(self Widget, key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey) bool
	KeyboardCharacterEvent(self Widget, codePoint rune) bool
	PreferredSize(self Widget, ctx *nanovgo.Context) (int, int)
	OnPerformLayout(self Widget, ctx *nanovgo.Context)
	Draw(ctx *nanovgo.Context)
	String() string
}

type WidgetImplement struct {
	parent                     Widget
	layout                     Layout
	theme                      *Theme
	x, y, w, h, fixedW, fixedH int
	visible, enabled           bool
	focused, mouseFocus        bool
	id                         string
	tooltip                    string
	fontSize                   int
	cursor                     Cursor
	children                   []Widget
}

func NewWidget(parent Widget) Widget {
	widget := &WidgetImplement{}
	InitWidget(widget, parent)
	return widget
}

// Parent() returns the parent widget
func (w *WidgetImplement) Parent() Widget {
	return w.parent
}

// SetParent() set the parent widget
func (w *WidgetImplement) SetParent(parent Widget) {
	w.parent = parent
}

// Layout() returns the used layout generator
func (w *WidgetImplement) Layout() Layout {
	return w.layout
}

// SetLayout() set the used layout generator
func (w *WidgetImplement) SetLayout(layout Layout) {
	w.layout = layout
}

// Theme() returns the theme used to draw this widget
func (w *WidgetImplement) Theme() *Theme {
	return w.theme
}

// SetTheme() set the theme used to draw this widget
func (w *WidgetImplement) SetTheme(theme *Theme) {
	w.theme = theme
}

// Position() returns the position relative to the parent widget
func (w *WidgetImplement) Position() (int, int) {
	return w.x, w.y
}

// SetPosition() set the position relative to the parent widget
func (w *WidgetImplement) SetPosition(x, y int) {
	w.x = x
	w.y = y
}

// AbsolutePosition() returns the absolute position on screen
func (w *WidgetImplement) AbsolutePosition() (int, int) {
	if w.parent != nil {
		x, y := w.parent.AbsolutePosition()
		return x + w.x, y + w.y
	}
	return w.x, w.y
}

// Size() returns the size of the widget
func (w *WidgetImplement) Size() (int, int) {
	return w.w, w.h
}

// SetSize() set the size of the widget
func (wg *WidgetImplement) SetSize(w, h int) {
	wg.w = w
	wg.h = h
}

// Width() returns the width of the widget
func (w *WidgetImplement) Width() int {
	return w.w
}

// SetWidth() set the width of the widget
func (wg *WidgetImplement) SetWidth(w int) {
	wg.w = w
}

// Height() returns the height of the widget
func (w *WidgetImplement) Height() int {
	return w.h
}

// SetHeight() set the height of the widget
func (w *WidgetImplement) SetHeight(h int) {
	w.h = h
}

// Return the fixed size (see SetFixedSize())
func (w *WidgetImplement) FixedSize() (int, int) {
	return w.fixedW, w.fixedH
}

// SetFixedSize() set the fixed size of this widget.
// If nonzero, components of the fixed size attribute override any values
// computed by a layout generator associated with this widget. Note that
// just setting the fixed size alone is not enough to actually change its
// size; this is done with a call to \ref SetSize or a call to PerformLayout()
// in the parent widget.
func (wg *WidgetImplement) SetFixedSize(w, h int) {
	wg.fixedW = w
	wg.fixedH = h
}

// FixedWidth() returns the fixed width (see SetFixedSize())
func (w *WidgetImplement) FixedWidth() int {
	return w.fixedW
}

// FixedHeight() returns the fixed height (see SetFixedSize())
func (w *WidgetImplement) FixedHeight() int {
	return w.fixedH
}

// SetFixedWidth() set the fixed width (see SetFixedSize())
func (wg *WidgetImplement) SetFixedWidth(w int) {
	wg.fixedW = w
}

// SetFixedSize() set the fixed height (see SetFixedSize())
func (w *WidgetImplement) SetFixedHeight(h int) {
	w.fixedH = h
}

// Visible() returns whether or not the widget is currently visible (assuming all parents are visible)
func (w *WidgetImplement) Visible() bool {
	return w.visible
}

// SetVisible() set whether or not the widget is currently visible (assuming all parents are visible)
func (w *WidgetImplement) SetVisible(v bool) {
	w.visible = v
}

// VisibleRecursive() checks if this widget is currently visible, taking parent widgets into account
func (w *WidgetImplement) VisibleRecursive() bool {
	if w.parent != nil {
		return w.Visible() && w.parent.VisibleRecursive()
	}
	return w.Visible()
}

// ChildCount() returns the number of child widgets
func (w *WidgetImplement) ChildCount() int {
	return len(w.children)
}

// Children() returns the list of child widgets of the current widget
func (w *WidgetImplement) Children() []Widget {
	return w.children
}

// AddChild() adds a child widget to the current widget
// This function almost never needs to be called by hand,
// since the constructor of \ref Widget automatically
// adds the current widget to its parent
func (w *WidgetImplement) AddChild(self, child Widget) {
	w.children = append(w.children, child)
	child.SetParent(self)
}

// RemoveChildByIndex() removes a child widget by index
func (w *WidgetImplement) RemoveChildByIndex(i int) {
	child := w.children[i]
	child.SetParent(nil)
	w.children, w.children[len(w.children)-1] = append(w.children[:i], w.children[i+1:]...), nil
}

// RemoveChild() removes a child widget by value
func (wg *WidgetImplement) RemoveChild(w Widget) {
	for i, child := range wg.children {
		if w == child {
			wg.RemoveChildByIndex(i)
			return
		}
	}
}

// Window() walks up the hierarchy and return the parent window
func (w *WidgetImplement) FindWindow() IWindow {
	parent := w.Parent()
	if parent == nil {
		panic("Widget:internal error (could not find parent window)")
	}
	return parent.FindWindow()
}

// SetID() associates this widget with an ID value (optional)
func (w *WidgetImplement) SetID(id string) {
	w.id = id
}

// ID() returns the ID value associated with this widget, if any
func (w *WidgetImplement) ID() string {
	return w.id
}

// Enabled() returns whether or not this widget is currently enabled
func (w *WidgetImplement) Enabled() bool {
	return w.enabled
}

/// SetEnabled() set whether or not this widget is currently enabled
func (w *WidgetImplement) SetEnabled(e bool) {
	w.enabled = e
}

// Focused() returns whether or not this widget is currently focused
func (w *WidgetImplement) Focused() bool {
	return w.focused
}

// SetFocused() set whether or not this widget is currently focused
func (w *WidgetImplement) SetFocused(f bool) {
	w.focused = f
}

// RequestFocus() requests the focus to be moved to this widget
func (w *WidgetImplement) RequestFocus(self Widget) {
	var widget Widget = self
	var parent Widget = self.Parent()
	for parent != nil {
		widget = parent
		parent = widget.Parent()
	}
	screen := widget.(*Screen)
	screen.UpdateFocus(self)
}

// Tooltip() returns tooltip string
func (w *WidgetImplement) Tooltip() string {
	return w.tooltip
}

// SetTooltip() set tooltip string
func (w *WidgetImplement) SetTooltip(s string) {
	w.tooltip = s
}

// FontSize() returns current font size. If not set the default of the current theme will be returned
func (w *WidgetImplement) FontSize() int {
	if w.fontSize > 0 {
		return w.fontSize
	}
	return w.theme.StandardFontSize
}

// SetFontSize() set the font size of this widget
func (w *WidgetImplement) SetFontSize(s int) {
	w.fontSize = s
}

// HasFontSize() return whether the font size is explicitly specified for this widget
func (w *WidgetImplement) HasFontSize() bool {
	return w.fontSize > 0
}

// Cursor() returns a pointer to the cursor of the widget
func (w *WidgetImplement) Cursor() Cursor {
	return w.cursor
}

// SetCursor() set the cursor of the widget
func (w *WidgetImplement) SetCursor(c Cursor) {
	w.cursor = c
}

// Contains() checks if the widget contains a certain position
func (w *WidgetImplement) Contains(x, y int) bool {
	return w.x <= x && w.y <= y && x <= w.x+w.w && y <= w.y+w.h
}

// FindWidget() determines the widget located at the given position value (recursive)
func (w *WidgetImplement) FindWidget(self Widget, x, y int) Widget {
	children := self.Children()
	for i := len(children) - 1; i > -1; i-- {
		child := children[i]
		if child.Visible() && child.Contains(x-w.x, y-w.y) {
			return child.FindWidget(child, x-w.x, y-w.y)
		}
	}
	if self.Contains(x, y) {
		return self
	}
	return nil
}

// MouseButtonEvent() handles a mouse button event (default implementation: propagate to children)
func (w *WidgetImplement) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	children := self.Children()
	for i := len(children) - 1; i > -1; i-- {
		child := children[i]
		if child.Visible() && child.Contains(x-w.x, y-w.y) && child.MouseButtonEvent(child, x-w.x, y-w.y, button, down, modifier) {
			return true
		}
	}
	if button == glfw.MouseButton1 && down && !w.focused {
		self.RequestFocus(self)
	}
	return false
}

// MouseMotionEvent() handles a mouse motion event (default implementation: propagate to children)
func (w *WidgetImplement) MouseMotionEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	children := self.Children()
	for i := len(children) - 1; i > -1; i-- {
		child := children[i]
		if !child.Visible() {
			continue
		}
		contained := child.Contains(x-w.x, y-w.y)
		prevContained := child.Contains(x-w.x-relX, y-w.y-relY)
		if contained != prevContained {
			child.MouseEnterEvent(child, x, y, contained)
		}
		if (contained || prevContained) && child.MouseMotionEvent(child, x-w.x, y-w.y, relX, relY, button, modifier) {
			return true
		}
	}
	return false
}

// MouseDragEvent() handles a mouse drag event (default implementation: do nothing)
func (w *WidgetImplement) MouseDragEvent(self Widget, x, y, relX, relY int, button int, modifier glfw.ModifierKey) bool {
	return false
}

// MouseEnterEvent() handles a mouse enter/leave event (default implementation: record this fact, but do nothing)
func (w *WidgetImplement) MouseEnterEvent(self Widget, x, y int, enter bool) bool {
	w.mouseFocus = enter
	return false
}

// ScrollEvent() handles a mouse scroll event (default implementation: propagate to children)
func (w *WidgetImplement) ScrollEvent(self Widget, x, y, relX, relY int) bool {
	children := self.Children()
	for i := len(children) - 1; i > -1; i-- {
		child := children[i]
		if !child.Visible() {
			continue
		}
		if child.Contains(x-w.x, y-w.y) && child.ScrollEvent(child, x-w.x, y-w.y, relX, relY) {
			return true
		}
	}
	return false
}

// FocusEvent() handles a focus change event (default implementation: record the focus status, but do nothing)
func (w *WidgetImplement) FocusEvent(self Widget, f bool) bool {
	w.focused = f
	return false
}

// KeyboardEvent() handles a keyboard event (default implementation: do nothing)
func (w *WidgetImplement) KeyboardEvent(self Widget, key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey) bool {
	return false
}

// KeyboardCharacterEvent() handles text input (UTF-32 format) (default implementation: do nothing)
func (w *WidgetImplement) KeyboardCharacterEvent(self Widget, codePoint rune) bool {
	return false
}

// PreferredSize() computes the preferred size of the widget
func (w *WidgetImplement) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	if w.layout != nil {
		return w.layout.PreferredSize(self, ctx)
	}
	return w.x, w.y
}

// PerformLayout() invokes the associated layout generator to properly place child widgets, if any
func (w *WidgetImplement) OnPerformLayout(self Widget, ctx *nanovgo.Context) {
	if w.layout != nil {
		w.layout.OnPerformLayout(self, ctx)
	} else {
		for _, child := range w.children {
			prefW, prefH := child.PreferredSize(child, ctx)
			fixW, fixH := child.FixedSize()
			w := toI(fixW > 0, fixW, prefW)
			h := toI(fixH > 0, fixH, prefH)
			child.SetSize(w, h)
			child.OnPerformLayout(child, ctx)
		}
	}
}

// Draw() draws the widget (and all child widgets)
func (w *WidgetImplement) Draw(ctx *nanovgo.Context) {
	if debug {
		ctx.SetStrokeWidth(1.0)
		ctx.BeginPath()
		ctx.Rect(float32(w.x)-0.5, float32(w.y)-0.5, float32(w.w)+1.0, float32(w.h)+1.0)
		ctx.SetStrokeColor(nanovgo.RGBA(255, 0, 0, 255))
		ctx.Stroke()
	}

	if len(w.children) == 0 {
		return
	}
	ctx.Translate(float32(w.x), float32(w.y))
	for _, child := range w.children {
		if child.Visible() {
			child.Draw(ctx)
		}
	}
	ctx.Translate(-float32(w.x), -float32(w.y))
}

func (w *WidgetImplement) String() string {
	return fmt.Sprintf("WidgetImplement [%d,%d-%d,%d]", w.x, w.y, w.w, w.h)
}
