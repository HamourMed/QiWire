package linuxinput

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	inputEventSize = 24
	evKey          = 0x01
)

type LinuxInputProvider struct {
	file *os.File
}

func NewLinuxInputProvider(path string) (*LinuxInputProvider, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &LinuxInputProvider{
		file: file,
	}, nil
}

func (p *LinuxInputProvider) ReadEvent() (KeyEvent, error) {
	buf := make([]byte, inputEventSize)

	for {
		_, err := io.ReadFull(p.file, buf)
		if err != nil {
			return KeyEvent{}, err
		}

		eventType := binary.LittleEndian.Uint16(buf[16:18])
		code := binary.LittleEndian.Uint16(buf[18:20])
		value := int32(binary.LittleEndian.Uint32(buf[20:24]))

		if eventType != evKey {
			continue
		}

		var action KeyAction

		switch value {
		case 0:
			action = KeyRelease
		case 1:
			action = KeyPress
		case 2:
			action = KeyRepeat
		default:
			return KeyEvent{}, fmt.Errorf("unknown key event value: %d", value)
		}

		return KeyEvent{
			Code:   code,
			Action: action,
		}, nil
	}
}

func (p *LinuxInputProvider) Close() error {
	return p.file.Close()
}
