package spx

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type Key = ebiten.Key

const (
	Key0            Key = ebiten.Key0
	Key1            Key = ebiten.Key1
	Key2            Key = ebiten.Key2
	Key3            Key = ebiten.Key3
	Key4            Key = ebiten.Key4
	Key5            Key = ebiten.Key5
	Key6            Key = ebiten.Key6
	Key7            Key = ebiten.Key7
	Key8            Key = ebiten.Key8
	Key9            Key = ebiten.Key9
	KeyA            Key = ebiten.KeyA
	KeyB            Key = ebiten.KeyB
	KeyC            Key = ebiten.KeyC
	KeyD            Key = ebiten.KeyD
	KeyE            Key = ebiten.KeyE
	KeyF            Key = ebiten.KeyF
	KeyG            Key = ebiten.KeyG
	KeyH            Key = ebiten.KeyH
	KeyI            Key = ebiten.KeyI
	KeyJ            Key = ebiten.KeyJ
	KeyK            Key = ebiten.KeyK
	KeyL            Key = ebiten.KeyL
	KeyM            Key = ebiten.KeyM
	KeyN            Key = ebiten.KeyN
	KeyO            Key = ebiten.KeyO
	KeyP            Key = ebiten.KeyP
	KeyQ            Key = ebiten.KeyQ
	KeyR            Key = ebiten.KeyR
	KeyS            Key = ebiten.KeyS
	KeyT            Key = ebiten.KeyT
	KeyU            Key = ebiten.KeyU
	KeyV            Key = ebiten.KeyV
	KeyW            Key = ebiten.KeyW
	KeyX            Key = ebiten.KeyX
	KeyY            Key = ebiten.KeyY
	KeyZ            Key = ebiten.KeyZ
	KeyApostrophe   Key = ebiten.KeyApostrophe
	KeyBackslash    Key = ebiten.KeyBackslash
	KeyBackspace    Key = ebiten.KeyBackspace
	KeyCapsLock     Key = ebiten.KeyCapsLock
	KeyComma        Key = ebiten.KeyComma
	KeyDelete       Key = ebiten.KeyDelete
	KeyDown         Key = ebiten.KeyDown
	KeyEnd          Key = ebiten.KeyEnd
	KeyEnter        Key = ebiten.KeyEnter
	KeyEqual        Key = ebiten.KeyEqual
	KeyEscape       Key = ebiten.KeyEscape
	KeyF1           Key = ebiten.KeyF1
	KeyF2           Key = ebiten.KeyF2
	KeyF3           Key = ebiten.KeyF3
	KeyF4           Key = ebiten.KeyF4
	KeyF5           Key = ebiten.KeyF5
	KeyF6           Key = ebiten.KeyF6
	KeyF7           Key = ebiten.KeyF7
	KeyF8           Key = ebiten.KeyF8
	KeyF9           Key = ebiten.KeyF9
	KeyF10          Key = ebiten.KeyF10
	KeyF11          Key = ebiten.KeyF11
	KeyF12          Key = ebiten.KeyF12
	KeyGraveAccent  Key = ebiten.KeyGraveAccent
	KeyHome         Key = ebiten.KeyHome
	KeyInsert       Key = ebiten.KeyInsert
	KeyKP0          Key = ebiten.KeyKP0
	KeyKP1          Key = ebiten.KeyKP1
	KeyKP2          Key = ebiten.KeyKP2
	KeyKP3          Key = ebiten.KeyKP3
	KeyKP4          Key = ebiten.KeyKP4
	KeyKP5          Key = ebiten.KeyKP5
	KeyKP6          Key = ebiten.KeyKP6
	KeyKP7          Key = ebiten.KeyKP7
	KeyKP8          Key = ebiten.KeyKP8
	KeyKP9          Key = ebiten.KeyKP9
	KeyKPDecimal    Key = ebiten.KeyKPDecimal
	KeyKPDivide     Key = ebiten.KeyKPDivide
	KeyKPEnter      Key = ebiten.KeyKPEnter
	KeyKPEqual      Key = ebiten.KeyKPEqual
	KeyKPMultiply   Key = ebiten.KeyKPMultiply
	KeyKPSubtract   Key = ebiten.KeyKPSubtract
	KeyLeft         Key = ebiten.KeyLeft
	KeyLeftBracket  Key = ebiten.KeyLeftBracket
	KeyMenu         Key = ebiten.KeyMenu
	KeyMinus        Key = ebiten.KeyMinus
	KeyNumLock      Key = ebiten.KeyNumLock
	KeyPageDown     Key = ebiten.KeyPageDown
	KeyPageUp       Key = ebiten.KeyPageUp
	KeyPause        Key = ebiten.KeyPause
	KeyPeriod       Key = ebiten.KeyPeriod
	KeyPrintScreen  Key = ebiten.KeyPrintScreen
	KeyRight        Key = ebiten.KeyRight
	KeyRightBracket Key = ebiten.KeyRightBracket
	KeyScrollLock   Key = ebiten.KeyScrollLock
	KeySemicolon    Key = ebiten.KeySemicolon
	KeySlash        Key = ebiten.KeySlash
	KeySpace        Key = ebiten.KeySpace
	KeyTab          Key = ebiten.KeyTab
	KeyUp           Key = ebiten.KeyUp
	KeyAlt          Key = ebiten.KeyAlt
	KeyControl      Key = ebiten.KeyControl
	KeyShift        Key = ebiten.KeyShift
	KeyMax          Key = ebiten.KeyMax
	KeyAny          Key = -1
)

// -------------------------------------------------------------------------------------

type event interface{}

type eventStart struct{}

type eventKeyDown struct {
	Key ebiten.Key
}

type eventKeyUp struct {
	Key ebiten.Key
}

type eventLeftButtonDown struct {
	X, Y int
}

type eventLeftButtonUp struct {
	X, Y int
}

type eventFirer interface {
	fireEvent(ev event)
}

type inputMgr struct {
	keyStates   map[ebiten.Key]int
	lbtnState   int
	keyDuration int
	firer       eventFirer
	startFlag   sync.Once
}

func (i *inputMgr) init(firer eventFirer, keyDuration int) {
	const (
		defaultKeyDuration = 15
	)
	if keyDuration == 0 {
		keyDuration = defaultKeyDuration
	}
	i.keyStates = make(map[ebiten.Key]int)
	i.lbtnState = mouseStateNone
	i.keyDuration = keyDuration
	i.firer = firer
}

func (i *inputMgr) reset() {
	i.startFlag = sync.Once{}
}

// -------------------------------------------------------------------------------------

const (
	mouseStateNone     = 0
	mouseStatePressing = 1
)

func (i *inputMgr) update() {
	i.startFlag.Do(func() {
		i.firer.fireEvent(&eventStart{})
	})
	i.updateKeyboard()
	i.updateMouse()
}

func (i *inputMgr) updateMouse() {
	switch i.lbtnState {
	case mouseStateNone:
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			i.lbtnState = mouseStatePressing
			i.firer.fireEvent(&eventLeftButtonDown{X: x, Y: y})
		}
	case mouseStatePressing:
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			i.lbtnState = mouseStateNone
			i.firer.fireEvent(&eventLeftButtonUp{X: x, Y: y})
		}
	default:
		panic("unknown mouse state")
	}
}

func (i *inputMgr) updateKeyboard() {
	keyDuration := i.keyDuration
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		if ebiten.IsKeyPressed(key) {
			n := i.keyStates[key]
			if n > 0 {
				if !isStateKey(key) {
					n--
				}
			}
			if n <= 0 {
				n = keyDuration
				i.firer.fireEvent(&eventKeyDown{Key: key})
			}
			i.keyStates[key] = n
		} else {
			if i.keyStates[key] > 0 {
				i.firer.fireEvent(&eventKeyUp{Key: key})
				i.keyStates[key] = 0
			}
		}
	}
}

func isStateKey(key Key) bool {
	switch key {
	case KeyAlt, KeyControl, KeyShift:
		return true
	}
	return false
}

func isKeyPressed(key Key) bool {
	if key == KeyAny {
		for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
			if ebiten.IsKeyPressed(key) {
				return true
			}
		}
		return false
	}
	return ebiten.IsKeyPressed(key)
}

func isMousePressed() bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

// -------------------------------------------------------------------------------------
