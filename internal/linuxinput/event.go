package linuxinput

const (
	evKey = 0x01

	keyEnter = 28
	keyA     = 30
	keyZ     = 44
	keySpace = 57

	maxKeyCode = 0x2ff
)

type KeyAction int

const (
	KeyRelease KeyAction = iota
	KeyPress
	KeyRepeat
)

type KeyEvent struct {
	Code   uint16
	Action KeyAction
}
