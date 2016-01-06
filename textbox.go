package nanogui

import (
	"github.com/shibukawa/glfw"
	"github.com/shibukawa/nanovgo"
	"regexp"
	"strconv"
)

type TextAlignment int

const (
	TextCenter TextAlignment = iota
	TextLeft
	TextRight
)

type TextBox struct {
	WidgetImplement

	fontFace            string
	editable            bool
	committed           bool
	value               string
	defaultValue        string
	yankValue           []rune
	alignment           TextAlignment
	units               string
	unitImage           int
	format              *regexp.Regexp
	callback            func(string) bool
	validFormat         bool
	valueTemp           []rune
	cursorPos           int
	selectionPos        int
	mousePos            [2]int
	mouseDownPos        [2]int
	mouseDragPos        [2]int
	mouseDownModifier   glfw.ModifierKey
	textOffset          float32
	lastClick           float32
	preeditText         []rune
	preeditBlocks       []int
	preeditFocusedBlock int
}

func NewTextBox(parent Widget, values ...string) *TextBox {
	var value string
	switch len(values) {
	case 0:
		value = "Untitled"
	case 1:
		value = values[0]
	default:
		panic("NewTextBox can accept only one extra parameter (value)")
	}

	textBox := &TextBox{}
	InitWidget(textBox, parent)
	textBox.init(value)
	return textBox
}

func (t *TextBox) init(value string) {
	t.committed = true
	t.value = value
	t.unitImage = -1
	t.validFormat = true
	t.valueTemp = []rune(value)
	t.cursorPos = -1
	t.selectionPos = -1
	t.mousePos = [2]int{-1, -1}
	t.mouseDownPos = [2]int{-1, -1}
	t.mouseDragPos = [2]int{-1, -1}
	t.fontSize = t.theme.TextBoxFontSize
}

func (t *TextBox) Editable() bool {
	return t.editable
}

func (t *TextBox) SetEditable(e bool) {
	t.editable = e
}

func (t *TextBox) Value() string {
	return t.value
}

func (t *TextBox) SetValue(value string) {
	t.value = value
}

func (t *TextBox) DefaultValue() string {
	return t.defaultValue
}

func (t *TextBox) SetDefaultValue(value string) {
	t.defaultValue = value
}

func (t *TextBox) Alignment() TextAlignment {
	return t.alignment
}

func (t *TextBox) SetAlignment(a TextAlignment) {
	t.alignment = a
}

func (t *TextBox) Units() string {
	return t.units
}

func (t *TextBox) SetUnits(units string) {
	t.units = units
}

func (t *TextBox) UnitImage() int {
	return t.unitImage
}

func (t *TextBox) SetUnitImage(img int) {
	t.unitImage = img
}

func (t *TextBox) Font() string {
	if t.fontFace == "" {
		return t.theme.FontNormal
	}
	return t.fontFace
}

func (t *TextBox) SetFont(fontFace string) {
	t.fontFace = fontFace
}

func (t *TextBox) Format() string {
	return t.format.String()
}

func (t *TextBox) SetFormat(format string) error {
	var err error
	t.format, err = regexp.Compile(format)
	return err
}

func (t *TextBox) SetCallback(callback func(string) bool) {
	t.callback = callback
}

func (t *TextBox) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	t.WidgetImplement.MouseButtonEvent(self, x, y, button, down, modifier)

	if t.editable && t.Focused() && button == glfw.MouseButton1 && len(t.preeditText) == 0 {
		if down {
			t.mouseDownPos = [2]int{x, y}
			t.mouseDownModifier = modifier
			time := GetTime()
			if time-t.lastClick < 0.25 {
				/* Double-click: select all text */
				t.selectionPos = 0
				t.cursorPos = len(t.valueTemp)
				t.mouseDownPos = [2]int{-1, 1}
			}
			t.lastClick = time
		} else {
			t.mouseDownPos = [2]int{-1, -1}
			t.mouseDragPos = [2]int{-1, -1}
		}
		return true
	}
	return false
}

func (t *TextBox) MouseMotionEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if t.editable && t.Focused() {
		t.mousePos = [2]int{x, y}
		return true
	}
	return false
}

func (t *TextBox) MouseDragEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if t.editable && t.Focused() {
		t.mouseDragPos = [2]int{x, y}
		return true
	}
	return false
}

func (t *TextBox) MouseEnterEvent(self Widget, x, y int, enter bool) bool {
	t.WidgetImplement.MouseEnterEvent(self, x, y, enter)
	return false
}

func (t *TextBox) FocusEvent(self Widget, focused bool) bool {
	t.WidgetImplement.FocusEvent(self, focused)
	backup := t.value

	if t.editable {
		if focused {
			t.valueTemp = []rune(t.value)
			t.committed = false
			t.cursorPos = 0
		} else {
			if t.validFormat {
				if len(t.valueTemp) == 0 {
					t.value = t.defaultValue
				} else {
					t.value = string(t.valueTemp)
				}
			}

			if t.callback != nil && !t.callback(t.value) {
				t.value = backup
			}

			t.validFormat = true
			t.committed = true
			t.cursorPos = -1
			t.selectionPos = -1
			t.textOffset = 0
		}
		t.validFormat = len(t.valueTemp) == 0 || t.checkFormat(string(t.valueTemp))
	}
	return true
}

func (t *TextBox) KeyboardEvent(self Widget, key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey) bool {
	if t.editable && t.Focused() {
		if (action == glfw.Press || action == glfw.Repeat) && len(t.preeditText) == 0 {
			switch DetectEditAction(key, modifier) {
			case EditActionMoveLeft:
				if modifier == glfw.ModShift {
					t.selectionPos = toI(t.selectionPos == -1, t.cursorPos, t.selectionPos)
				} else {
					t.selectionPos = -1
				}
				if t.cursorPos > 0 {
					t.cursorPos--
				}
			case EditActionMoveRight:
				if modifier == glfw.ModShift {
					t.selectionPos = toI(t.selectionPos == -1, t.cursorPos, t.selectionPos)
				} else {
					t.selectionPos = -1
				}
				if t.cursorPos < len(t.valueTemp) {
					t.cursorPos++
				}
			case EditActionMoveLineTop:
				if modifier == glfw.ModShift {
					t.selectionPos = toI(t.selectionPos == -1, t.cursorPos, t.selectionPos)
				} else {
					t.selectionPos = -1
				}
				t.cursorPos = 0
			case EditActionMoveLineEnd:
				if modifier == glfw.ModShift {
					t.selectionPos = toI(t.selectionPos == -1, t.cursorPos, t.selectionPos)
				} else {
					t.selectionPos = -1
				}
				t.cursorPos = len(t.valueTemp)
			case EditActionBackspace:
				if !t.DeleteSelection() {
					if t.cursorPos > 0 {
						t.valueTemp = append(t.valueTemp[:t.cursorPos-1], t.valueTemp[t.cursorPos:]...)
						t.cursorPos--
					}
				}
			case EditActionDelete:
				if !t.DeleteSelection() {
					if t.cursorPos < len(t.valueTemp) {
						t.valueTemp = append(t.valueTemp[:t.cursorPos], t.valueTemp[t.cursorPos+1:]...)
					}
				}
			case EditActionCutUntilLineEnd:
				t.yankValue = t.valueTemp[t.cursorPos:]
				t.valueTemp = t.valueTemp[:t.cursorPos]
			case EditActionYank:
				t.valueTemp = append(t.valueTemp[:t.cursorPos], append(t.yankValue, t.valueTemp[t.cursorPos:]...)...)
			case EditActionDeleteLeftWord:
				panic("not implemented")
			case EditActionEnter:
				if !t.committed {
					t.FocusEvent(t, false)
				}
			case EditActionSelectAll:
				t.cursorPos = len(t.valueTemp)
				t.selectionPos = 0
			case EditActionCopy:
				t.CopySelection()
			case EditActionCut:
				t.CopySelection()
				t.DeleteSelection()
			case EditActionPaste:
				t.DeleteSelection()
				t.PasteFromClipboard()
			}
			t.validFormat = len(t.valueTemp) == 0 || t.checkFormat(string(t.valueTemp))
		}
		return true
	}
	return false
}

func (t *TextBox) KeyboardCharacterEvent(self Widget, codePoint rune) bool {
	if t.editable && t.Focused() {
		t.DeleteSelection()
		t.valueTemp = append(t.valueTemp[:t.cursorPos], append([]rune{codePoint}, t.valueTemp[t.cursorPos:]...)...)
		t.cursorPos++
		t.validFormat = len(t.valueTemp) == 0 || t.checkFormat(string(t.valueTemp))
		t.preeditText = nil
		return true
	}
	return false
}

func (t *TextBox) IMEPreeditEvent(self Widget, text []rune, blocks []int, focusedBlock int) bool {
	t.preeditText = text
	t.preeditBlocks = blocks
	t.preeditFocusedBlock = focusedBlock
	return true
}

func (t *TextBox) IMEStatusEvent(self Widget) bool {
	if len(t.preeditText) != 0 {
		t.valueTemp = append(append(t.valueTemp[:t.cursorPos], t.preeditText...), t.valueTemp[t.cursorPos:]...)
		t.cursorPos += len(t.preeditText)
		t.preeditText = nil
	}
	return true
}

func (t *TextBox) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	sizeH := float32(t.FontSize()) * 1.4

	var unitWidth, textWidth float32
	ctx.SetFontSize(float32(t.FontSize()))
	if t.unitImage > 0 {
		w, h, _ := ctx.ImageSize(t.unitImage)
		unitHeight := sizeH * 0.4
		unitWidth = float32(w) * unitHeight / float32(h)
	} else if t.units != "" {
		unitWidth, _ = ctx.TextBounds(0, 0, t.units)
	}

	textWidth, _ = ctx.TextBounds(0, 0, string(t.editingText()))
	sizeW := sizeH + textWidth + unitWidth
	return int(sizeW), int(sizeH)
}

func (t *TextBox) Draw(self Widget, ctx *nanovgo.Context) {
	t.WidgetImplement.Draw(self, ctx)

	x := float32(t.x)
	y := float32(t.y)
	w := float32(t.w)
	h := float32(t.h)

	bg := nanovgo.BoxGradient(x+1, y+2, w-2, h-2, 3, 4, nanovgo.MONO(255, 32), nanovgo.MONO(32, 32))
	fg1 := nanovgo.BoxGradient(x+1, y+2, w-2, h-2, 3, 4, nanovgo.MONO(150, 32), nanovgo.MONO(32, 32))
	fg2 := nanovgo.BoxGradient(x+1, y+2, w-2, h-2, 3, 4, nanovgo.RGBA(255, 0, 0, 100), nanovgo.RGBA(255, 0, 0, 50))

	ctx.BeginPath()
	ctx.RoundedRect(x+1, y+2, w-2, h-2, 3)
	if t.editable && t.Focused() {
		if t.validFormat {
			ctx.SetFillPaint(fg1)
		} else {
			ctx.SetFillPaint(fg2)
		}
	} else {
		ctx.SetFillPaint(bg)
	}

	ctx.Fill()

	ctx.BeginPath()
	ctx.RoundedRect(x+0.5, y+0.5, w-1, h-1, 2.5)
	ctx.SetStrokeColor(nanovgo.MONO(0, 48))
	ctx.Stroke()

	ctx.SetFontSize(float32(t.FontSize()))
	ctx.SetFontFace(t.Font())
	drawPosX := x
	drawPosY := y + h*0.5 + 1

	xSpacing := h * 0.3
	var unitWidth float32

	if t.unitImage > 0 {
		iw, ih, _ := ctx.ImageSize(t.unitImage)
		unitHeight := float32(ih) * 0.4
		unitWidth = float32(iw) * unitHeight / float32(h)
		imgPaint := nanovgo.ImagePattern(x+w-xSpacing-unitWidth, drawPosY-unitHeight*0.5,
			unitWidth, unitHeight, 0, t.unitImage, toF(t.enabled, 0.7, 0.35))
		ctx.BeginPath()
		ctx.Rect(x+w-xSpacing-unitWidth, drawPosY-unitHeight*0.5, unitWidth, unitHeight)
		ctx.SetFillPaint(imgPaint)
		ctx.Fill()
		unitWidth += 2
	} else if t.units != "" {
		unitWidth, _ = ctx.TextBounds(0, 0, t.units)
		ctx.SetFillColor(nanovgo.MONO(255, toB(t.enabled, 64, 32)))
		ctx.SetTextAlign(nanovgo.AlignRight | nanovgo.AlignMiddle)
		ctx.Text(x+w-xSpacing, drawPosY, t.units)
	}

	switch t.alignment {
	case TextLeft:
		ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignMiddle)
		drawPosX += xSpacing
	case TextRight:
		ctx.SetTextAlign(nanovgo.AlignRight | nanovgo.AlignMiddle)
		drawPosX += w - unitWidth - xSpacing
	case TextCenter:
		ctx.SetTextAlign(nanovgo.AlignCenter | nanovgo.AlignMiddle)
		drawPosX += w * 0.5
	}
	if t.enabled {
		ctx.SetFillColor(t.theme.TextColor)
	} else {
		ctx.SetFillColor(t.theme.DisabledTextColor)
	}
	// clip visible text area
	clipX := x + xSpacing - 1
	clipY := y + 1.0
	clipWidth := w - unitWidth - 2.0*xSpacing + 2.0
	clipHeight := h - 3.0
	ctx.Scissor(clipX, clipY, clipWidth, clipHeight)
	oldDrawPosX := drawPosX
	drawPosX += t.textOffset

	if t.committed {
		ctx.Text(drawPosX, drawPosY, t.value)
	} else {
		text := t.editingText()
		textString := string(text)
		_, bounds := ctx.TextBounds(drawPosX, drawPosY, textString)
		lineH := bounds[3] - bounds[1]
		// find cursor positions
		glyphs := ctx.TextGlyphPositionsRune(drawPosX, drawPosY, text)
		t.updateCursor(ctx, bounds[2], glyphs)

		// compute text offset
		prevCPos := toI(t.cursorPos > 0, t.cursorPos-1, 0)
		nextCPos := toI(t.cursorPos < len(glyphs), t.cursorPos+1, len(glyphs))
		prevCX := t.textIndex2Position(prevCPos, bounds[2], glyphs)
		nextCX := t.textIndex2Position(nextCPos, bounds[2], glyphs)

		if nextCX > clipX+clipWidth {
			t.textOffset -= nextCX - (clipX + clipWidth) + 1.0
		}
		if prevCX < clipX {
			t.textOffset += clipX - prevCX + 1.0
		}
		drawPosX = oldDrawPosX + t.textOffset

		// draw text with offset
		ctx.TextRune(drawPosX, drawPosY, text)
		_, bounds = ctx.TextBounds(drawPosX, drawPosY, textString)

		// recompute cursor position
		glyphs = ctx.TextGlyphPositionsRune(drawPosX, drawPosY, text)

		var caretX float32 = -1
		if len(t.preeditText) != 0 {
			// draw preedit text
			caretX = t.textIndex2Position(t.cursorPos+len(t.preeditText), bounds[2], glyphs)

			offsetIndex := t.cursorPos
			offsetX := t.textIndex2Position(t.cursorPos, bounds[2], glyphs)
			ctx.SetStrokeColor(nanovgo.MONO(255, 160))
			ctx.SetFillColor(nanovgo.MONO(255, 80))
			ctx.SetStrokeWidth(2.0)
			for i, blockLength := range t.preeditBlocks {
				nextOffsetIndex := offsetIndex + blockLength
				nextOffsetX := t.textIndex2Position(nextOffsetIndex, bounds[2], glyphs)
				if i != t.preeditFocusedBlock {
					ctx.BeginPath()
					ctx.MoveTo(offsetX+2, drawPosY+lineH*0.5-1)
					ctx.LineTo(nextOffsetX-2, drawPosY+lineH*0.5-1)
					ctx.Stroke()
				} else {
					ctx.BeginPath()
					ctx.Rect(offsetX, drawPosY-lineH*0.5, nextOffsetX-offsetX, lineH)
					ctx.Fill()
				}
				offsetIndex = nextOffsetIndex
				offsetX = nextOffsetX
			}
			screen := t.FindWindow().Parent().(*Screen)
			oldCurX, oldCurY, oldCurH := screen.PreeditCursorPos()
			absX, absY := t.Parent().AbsolutePosition()
			newCurX := int(caretX) + absX
			newCurY := int(drawPosY+lineH*0.5) + absY
			newCurH := int(lineH)
			if oldCurX != newCurX || oldCurY != newCurY || oldCurH != newCurH {
				screen.SetPreeditCursorPos(newCurX, newCurY, newCurH)
			}
		} else if t.cursorPos > -1 {
			// regular cursor and selection area
			caretX = t.textIndex2Position(t.cursorPos, bounds[2], glyphs)

			if t.selectionPos > -1 {
				caretX2 := caretX
				selX := t.textIndex2Position(t.selectionPos, bounds[2], glyphs)

				if caretX2 > selX {
					selX, caretX2 = caretX2, selX
				}

				// draw selection
				ctx.BeginPath()
				ctx.SetFillColor(nanovgo.MONO(255, 80))
				ctx.Rect(caretX2, drawPosY-lineH*0.5, selX-caretX2, lineH)
				ctx.Fill()
			}
		}
		if caretX > 0 {
			// draw cursor
			ctx.BeginPath()
			ctx.MoveTo(caretX, drawPosY-lineH*0.5)
			ctx.LineTo(caretX, drawPosY+lineH*0.5)
			ctx.SetStrokeColor(nanovgo.RGBA(255, 192, 0, 255))
			ctx.SetStrokeWidth(1.0)
			ctx.Stroke()
		}
	}
	ctx.ResetScissor()
}

func (t *TextBox) checkFormat(input string) bool {
	if t.format == nil {
		return true
	}
	return t.format.MatchString(input)
}

func (t *TextBox) CopySelection() bool {
	sc := t.FindWindow().Parent().(*Screen)
	if t.selectionPos > -1 {
		begin := t.cursorPos
		end := t.selectionPos

		if begin > end {
			begin, end = end, begin
		}
		sc.GLFWWindow().SetClipboardString(string(t.valueTemp[begin:end]))
	}
	return false
}

func (t *TextBox) PasteFromClipboard() {
	sc := t.FindWindow().Parent().(*Screen)
	str, _ := sc.GLFWWindow().GetClipboardString()
	runes := []rune(str)
	t.valueTemp = append(t.valueTemp[:t.cursorPos], append(runes, t.valueTemp[t.cursorPos:]...)...)
	t.cursorPos += len(runes)
}

func (t *TextBox) DeleteSelection() bool {
	if t.selectionPos > -1 {
		begin := t.cursorPos
		end := t.selectionPos

		if begin > end {
			begin, end = end, begin
		}
		t.valueTemp = append(t.valueTemp[:begin], t.valueTemp[end:]...)
		t.cursorPos = begin
		t.selectionPos = -1
		return true
	}
	return false
}

func (t *TextBox) updateCursor(ctx *nanovgo.Context, lastX float32, glyphs []nanovgo.GlyphPosition) {
	if t.mouseDownPos[0] != -1 {
		if t.mouseDownModifier == glfw.ModShift {
			if t.selectionPos == -1 {
				t.selectionPos = t.cursorPos
			}
		} else {
			t.selectionPos = -1
		}
		t.cursorPos = t.position2CursorIndex(float32(t.mouseDownPos[0]), lastX, glyphs)
		t.mouseDownPos = [2]int{-1, -1}
	} else if t.mouseDragPos[0] != -1 {
		if t.selectionPos == -1 {
			t.selectionPos = t.cursorPos
		}
		t.cursorPos = t.position2CursorIndex(float32(t.mouseDragPos[0]), lastX, glyphs)
	} else {
		// set cursor to last character
		if t.cursorPos == -2 {
			t.cursorPos = len(glyphs)
		}
	}

	if t.cursorPos == t.selectionPos {
		t.selectionPos = -1
	}
}

func (t *TextBox) textIndex2Position(index int, lastX float32, glyphs []nanovgo.GlyphPosition) float32 {
	if index == len(glyphs) {
		return lastX
	}
	return glyphs[index].X
}

func (t *TextBox) position2CursorIndex(posX, lastX float32, glyphs []nanovgo.GlyphPosition) int {
	cursorIndex := 0
	if len(glyphs) == 0 {
		return 0
	}
	caretX := glyphs[0].X
	for j := 1; j < len(glyphs); j++ {
		glyph := &glyphs[j]
		if absF(caretX-posX) > absF(glyph.X-posX) {
			cursorIndex = j
			caretX = glyph.X
		}
	}
	if absF(caretX-posX) > absF(lastX-posX) {
		return len(glyphs)
	}
	return cursorIndex
}

func (t *TextBox) editingText() []rune {
	if len(t.preeditText) == 0 {
		return t.valueTemp
	}
	result := make([]rune, 0, len(t.valueTemp)+len(t.preeditText))
	result = append(append(append(result, t.valueTemp[:t.cursorPos]...), t.preeditText...), t.valueTemp[t.cursorPos:]...)
	return result
}

func (t *TextBox) String() string {
	return t.StringHelper("TextBox", t.value)
}

type IntBox struct {
	TextBox

	callback func(int)
}

func NewIntBox(parent Widget, signed bool, values ...int) *IntBox {
	var value int
	switch len(values) {
	case 0:
		value = 0
	case 1:
		value = values[0]
	default:
		panic("NewIntBox can accept only one extra parameter (value)")
	}

	intBox := &IntBox{}
	InitWidget(intBox, parent)
	intBox.init("")
	if signed {
		intBox.SetFormat(`^[-]?[0-9]*$`)
	} else {
		intBox.SetFormat(`^[0-9]*$`)
	}
	intBox.SetValue(value)
	return intBox
}

func (t *IntBox) Value() int {
	v, _ := strconv.ParseInt(t.value, 10, 64)
	return int(v)
}

func (t *IntBox) SetValue(value int) {
	t.value = strconv.FormatInt(int64(value), 10)
}

func (i *IntBox) DefaultValue() int {
	v, _ := strconv.ParseInt(i.defaultValue, 10, 64)
	return int(v)
}

func (i *IntBox) SetDefaultValue(value int) {
	i.defaultValue = strconv.FormatInt(int64(value), 10)
}

func (i *IntBox) SetCallback(callback func(int)) {
	i.callback = callback
}

func (i *IntBox) String() string {
	return i.StringHelper("IntBox", i.value)
}

type FloatBox struct {
	TextBox

	callback func(float64)
}

func NewFloatBox(parent Widget, values ...float64) *FloatBox {
	var value float64
	switch len(values) {
	case 0:
		value = 0.0
	case 1:
		value = values[0]
	default:
		panic("NewFloatBox can accept only one extra parameter (value)")
	}

	floatBox := &FloatBox{}
	InitWidget(floatBox, parent)
	floatBox.init("")
	floatBox.SetFormat(`^[-]?[0-9]*\.?[0-9]+$`)
	floatBox.SetValue(value)

	return floatBox
}

func (t *FloatBox) Value() float64 {
	v, _ := strconv.ParseFloat(string(t.value), 64)
	return v
}

func (t *FloatBox) SetValue(value float64) {
	// todo: remove trailing zero
	t.value = strconv.FormatFloat(value, 'f', 5, 64)
}

func (t *FloatBox) DefaultValue() float64 {
	v, _ := strconv.ParseFloat(string(t.defaultValue), 64)
	return v
}

func (t *FloatBox) SetDefaultValue(value float64) {
	// todo: remove trailing zero
	t.defaultValue = strconv.FormatFloat(value, 'f', 5, 64)
}

func (f *FloatBox) SetCallback(callback func(float64)) {
	f.callback = callback
}

func (f *FloatBox) String() string {
	return f.StringHelper("FloatBox", f.value)
}
