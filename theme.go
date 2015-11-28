package nanogui

import (
	"github.com/shibukawa/nanovgo"
)

type Theme struct {
	StandardFontSize     int
	ButtonFontSize       int
	TextBoxFontSize      int
	WindowCornerRadius   int
	WindowHeaderHeight   int
	WindowDropShadowSize int
	ButtonCornerRadius   int

	DropShadow        nanovgo.Color
	Transparent       nanovgo.Color
	BorderDark        nanovgo.Color
	BorderLight       nanovgo.Color
	BorderMedium      nanovgo.Color
	TextColor         nanovgo.Color
	DisabledTextColor nanovgo.Color
	TextColorShadow   nanovgo.Color
	IconColor         nanovgo.Color

	ButtonGradientTopFocused   nanovgo.Color
	ButtonGradientBotFocused   nanovgo.Color
	ButtonGradientTopUnfocused nanovgo.Color
	ButtonGradientBotUnfocused nanovgo.Color
	ButtonGradientTopPushed    nanovgo.Color
	ButtonGradientBotPushed    nanovgo.Color

	/* Window-related */
	WindowFillUnfocused  nanovgo.Color
	WindowFillFocused    nanovgo.Color
	WindowTitleUnfocused nanovgo.Color
	WindowTitleFocused   nanovgo.Color

	WindowHeaderGradientTop nanovgo.Color
	WindowHeaderGradientBot nanovgo.Color
	WindowHeaderSepTop      nanovgo.Color
	WindowHeaderSepBot      nanovgo.Color

	WindowPopup            nanovgo.Color
	WindowPopupTransparent nanovgo.Color

	FontNormal string
	FontBold   string
	FontIcons  string
}

func NewStandardTheme(ctx *nanovgo.Context) *Theme {
	ctx.CreateFontFromMemory("sans", MustAsset("fonts/Roboto-Regular.ttf"), 0)
	ctx.CreateFontFromMemory("sans-bold", MustAsset("fonts/Roboto-Bold.ttf"), 0)
	ctx.CreateFontFromMemory("icons", MustAsset("fonts/MaterialIcons-Regular.ttf"), 0)
	return &Theme{
		StandardFontSize:     16,
		ButtonFontSize:       20,
		TextBoxFontSize:      20,
		WindowCornerRadius:   2,
		WindowHeaderHeight:   30,
		WindowDropShadowSize: 10,
		ButtonCornerRadius:   2,

		DropShadow:        nanovgo.MONO(0, 128),
		Transparent:       nanovgo.MONO(0, 0),
		BorderDark:        nanovgo.MONO(29, 255),
		BorderLight:       nanovgo.MONO(92, 255),
		BorderMedium:      nanovgo.MONO(35, 255),
		TextColor:         nanovgo.MONO(255, 160),
		DisabledTextColor: nanovgo.MONO(255, 80),
		TextColorShadow:   nanovgo.MONO(0, 160),
		IconColor:         nanovgo.MONO(255, 160),

		ButtonGradientTopFocused:   nanovgo.MONO(64, 255),
		ButtonGradientBotFocused:   nanovgo.MONO(48, 255),
		ButtonGradientTopUnfocused: nanovgo.MONO(74, 255),
		ButtonGradientBotUnfocused: nanovgo.MONO(58, 255),
		ButtonGradientTopPushed:    nanovgo.MONO(41, 255),
		ButtonGradientBotPushed:    nanovgo.MONO(29, 255),

		WindowFillUnfocused:  nanovgo.MONO(43, 230),
		WindowFillFocused:    nanovgo.MONO(45, 230),
		WindowTitleUnfocused: nanovgo.MONO(220, 160),
		WindowTitleFocused:   nanovgo.MONO(255, 190),

		WindowHeaderGradientTop: nanovgo.MONO(74, 255),
		WindowHeaderGradientBot: nanovgo.MONO(58, 255),
		WindowHeaderSepTop:      nanovgo.MONO(92, 255),
		WindowHeaderSepBot:      nanovgo.MONO(29, 255),

		WindowPopup:            nanovgo.MONO(50, 255),
		WindowPopupTransparent: nanovgo.MONO(50, 0),

		FontNormal: "sans",
		FontBold:   "sans-bold",
		FontIcons:  "icons",
	}
}
