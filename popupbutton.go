package nanogui

import (
	"github.com/shibukawa/nanovgo"
	"runtime"
)

type PopupButton struct {
	Button
	chevronIcon Icon
	popup       *Popup
}

func finalizePopupButton(button *PopupButton) {
	if button.popup != nil {
		parent := button.popup.Parent()
		parent.RemoveChild(button.popup)
		button.popup = nil
	}
}

func NewPopupButton(parent Widget, captions ...string) *PopupButton {
	var caption string
	switch len(captions) {
	case 0:
		caption = "Untitled"
	case 1:
		caption = captions[0]
	default:
		panic("NewPopupButton can accept only one extra parameter (caption)")
	}

	button := &PopupButton{
		chevronIcon: IconRightOpen,
	}
	button.SetCaption(caption)
	button.SetIconPosition(ButtonIconLeftCentered)
	button.SetFlags(ToggleButtonType | PopupButtonType)

	parentWindow := parent.FindWindow()
	button.popup = NewPopup(parentWindow.Parent(), parentWindow)
	button.popup.SetSize(320, 250)

	InitWidget(button, parent)

	runtime.SetFinalizer(button, finalizePopupButton)

	return button
}

func (p *PopupButton) ChevronIcon() Icon {
	return p.chevronIcon
}

func (p *PopupButton) SetChevronIcon(icon Icon) {
	p.chevronIcon = icon
}

func (p *PopupButton) Popup() Widget {
	return p.popup.panel
}

func (p *PopupButton) Draw(ctx *nanovgo.Context) {
	if !p.enabled && p.pushed {
		p.pushed = false
	}
	p.popup.SetVisible(p.pushed)
	p.Button.Draw(ctx)
	if p.chevronIcon != 0 {
		ctx.SetFillColor(p.TextColor())
		ctx.SetFontSize(float32(p.FontSize()))
		ctx.SetFontFace(p.theme.FontIcons)
		ctx.SetTextAlign(nanovgo.AlignMiddle | nanovgo.AlignLeft)
		fontString := string([]rune{rune(p.chevronIcon)})
		iw, _ := ctx.TextBounds(0, 0, fontString)
		px, py := p.Position()
		w, h := p.Size()
		ix := px + w - int(iw) - 8
		iy := py + h/2 - 1
		ctx.Text(float32(ix), float32(iy), fontString)
	}
}

func (p *PopupButton) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	w, h := p.Button.PreferredSize(self, ctx)
	return w + 15, h
}

func (p *PopupButton) OnPerformLayout(self Widget, ctx *nanovgo.Context) {
	p.Button.WidgetImplement.OnPerformLayout(self, ctx)
	parentWindow := self.FindWindow()
	x := parentWindow.Width() + 15
	_, ay := p.AbsolutePosition()
	_, py := parentWindow.Position()
	y := ay - py + p.Height()/2
	p.popup.SetAnchorPosition(x, y)
}

func (p *PopupButton) String() string {
	return p.StringHelper("PopupButton", p.caption)
}
