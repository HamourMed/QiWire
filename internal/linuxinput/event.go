package linuxinput

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
