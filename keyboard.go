package nanogui

import (
	"github.com/goxjs/glfw"
	"runtime"
)

type EditAction int

const (
	EditActionNone EditAction = iota
	EditActionMoveLeft
	EditActionMoveRight
	EditActionMoveLineTop
	EditActionMoveLineEnd
	EditActionBackspace
	EditActionDelete
	EditActionCutUntilLineEnd
	EditActionYank
	EditActionDeleteLeftWord
	EditActionEnter
	EditActionSelectAll
	EditActionCopy
	EditActionCut
	EditActionPaste
)

func DetectEditAction(key glfw.Key, modifier glfw.ModifierKey) EditAction {
	isMac := runtime.GOOS == "darwin"
	switch key {
	case glfw.KeyLeft:
		return EditActionMoveLeft
	case glfw.KeyB:
		if modifier == glfw.ModControl {
			return EditActionMoveLeft
		}
	case glfw.KeyRight:
		return EditActionMoveRight
	case glfw.KeyF:
		if modifier == glfw.ModControl {
			return EditActionMoveRight
		}
	case glfw.KeyHome:
		return EditActionMoveLineTop
	case glfw.KeyA:
		if modifier == glfw.ModControl {
			if isMac {
				return EditActionMoveLineTop
			} else {
				return EditActionSelectAll
			}
		} else if modifier == glfw.ModSuper {
			return EditActionSelectAll
		}
	case glfw.KeyEnd:
		return EditActionMoveLineEnd
	case glfw.KeyE:
		if modifier == glfw.ModControl {
			return EditActionMoveLineEnd
		}
	case glfw.KeyBackspace:
		return EditActionBackspace
	case glfw.KeyH:
		if modifier == glfw.ModControl {
			return EditActionBackspace
		}
	case glfw.KeyDelete:
		if modifier == glfw.ModAlt {
			return EditActionDeleteLeftWord
		}
		return EditActionDelete
	case glfw.KeyD:
		if modifier == glfw.ModControl {
			return EditActionDelete
		}
	case glfw.KeyK:
		if modifier == glfw.ModControl {
			return EditActionCutUntilLineEnd
		}
	case glfw.KeyY:
		if modifier == glfw.ModControl {
			return EditActionYank
		}
	case glfw.KeyEnter:
		return EditActionEnter
	case glfw.KeyC:
		if (!isMac && modifier == glfw.ModControl) || (isMac && modifier == glfw.ModSuper) {
			return EditActionCopy
		}
	case glfw.KeyX:
		if (!isMac && modifier == glfw.ModControl) || (isMac && modifier == glfw.ModSuper) {
			return EditActionCut
		}
	case glfw.KeyV:
		if (!isMac && modifier == glfw.ModControl) || (isMac && modifier == glfw.ModSuper) {
			return EditActionPaste
		}
	}
	return EditActionNone
}
