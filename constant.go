package nanogui

type Cursor int

const (
	Arrow Cursor = iota
	IBeam
	Crosshair
	Hand
	HResize
	VResize
	CursorCount
)
