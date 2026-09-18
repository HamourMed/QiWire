package xkb

/*
#cgo pkg-config: xkbcommon

#include <stdlib.h>
#include <xkbcommon/xkbcommon.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"qiwire/internal/linuxinput"
)

type XKBTranslator struct {
	context *C.struct_xkb_context
	keymap  *C.struct_xkb_keymap
	state   *C.struct_xkb_state
}

func NewXKBTranslator(layout string) (*XKBTranslator, error) {
	ctx := C.xkb_context_new(C.XKB_CONTEXT_NO_FLAGS)
	if ctx == nil {
		return nil, fmt.Errorf("failed to create XKB context")
	}

	cLayout := C.CString(layout)
	defer C.free(unsafe.Pointer(cLayout))

	names := C.struct_xkb_rule_names{
		rules:   nil,
		model:   nil,
		layout:  cLayout,
		variant: nil,
		options: nil,
	}

	keymap := C.xkb_keymap_new_from_names(
		ctx,
		&names,
		C.XKB_KEYMAP_COMPILE_NO_FLAGS,
	)
	if keymap == nil {
		C.xkb_context_unref(ctx)
		return nil, fmt.Errorf("failed to create XKB keymap for layout %q", layout)
	}

	state := C.xkb_state_new(keymap)
	if state == nil {
		C.xkb_keymap_unref(keymap)
		C.xkb_context_unref(ctx)
		return nil, fmt.Errorf("failed to create XKB state")
	}

	return &XKBTranslator{
		context: ctx,
		keymap:  keymap,
		state:   state,
	}, nil
}

func (t *XKBTranslator) Translate(event linuxinput.KeyEvent) string {
	// Linux evdev keycodes use the common XKB offset of +8.
	keycode := C.xkb_keycode_t(event.Code + 8)

	switch event.Action {
	case linuxinput.KeyPress:
		text := t.keyText(keycode)

		C.xkb_state_update_key(
			t.state,
			keycode,
			C.XKB_KEY_DOWN,
		)

		return text

	case linuxinput.KeyRepeat:
		// The key is already considered down.
		return t.keyText(keycode)

	case linuxinput.KeyRelease:
		C.xkb_state_update_key(
			t.state,
			keycode,
			C.XKB_KEY_UP,
		)

		return ""

	default:
		return ""
	}
}

func (t *XKBTranslator) keyText(keycode C.xkb_keycode_t) string {
	keysym := C.xkb_state_key_get_one_sym(t.state, keycode)

	switch keysym {
	case C.XKB_KEY_Return, C.XKB_KEY_KP_Enter:
		return "\n"

	case C.XKB_KEY_BackSpace:
		return "\b"
	}

	size := C.xkb_state_key_get_utf8(
		t.state,
		keycode,
		nil,
		0,
	)

	if size <= 0 {
		return ""
	}

	buf := C.malloc(C.size_t(size + 1))
	if buf == nil {
		return ""
	}
	defer C.free(buf)

	C.xkb_state_key_get_utf8(
		t.state,
		keycode,
		(*C.char)(buf),
		C.size_t(size+1),
	)

	return C.GoString((*C.char)(buf))
}

func (t *XKBTranslator) Close() {
	if t.state != nil {
		C.xkb_state_unref(t.state)
		t.state = nil
	}

	if t.keymap != nil {
		C.xkb_keymap_unref(t.keymap)
		t.keymap = nil
	}

	if t.context != nil {
		C.xkb_context_unref(t.context)
		t.context = nil
	}
}
