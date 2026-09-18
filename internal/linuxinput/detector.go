package linuxinput

import (
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/unix"
)

type KeyboardDevice struct {
	Path string
	Name string
}

func DetectKeyboards() ([]KeyboardDevice, error) {
	paths, err := filepath.Glob("/dev/input/event*")
	if err != nil {
		return nil, err
	}

	var keyboards []KeyboardDevice

	for _, path := range paths {
		device, ok, err := inspectDevice(path)
		if err != nil {
			continue
		}

		if ok {
			keyboards = append(keyboards, device)
		}
	}

	return keyboards, nil
}

func inspectDevice(path string) (KeyboardDevice, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return KeyboardDevice{}, false, err
	}
	defer file.Close()

	fd := int(file.Fd())

	name, err := deviceName(fd)
	if err != nil {
		name = "unknown"
	}

	isKeyboard, err := hasKeyboardCapabilities(fd)
	if err != nil {
		return KeyboardDevice{}, false, err
	}

	if !isKeyboard {
		return KeyboardDevice{}, false, nil
	}

	return KeyboardDevice{
		Path: path,
		Name: name,
	}, true, nil
}

func hasKeyboardCapabilities(fd int) (bool, error) {
	const bitsPerByte = 8

	buf := make([]byte, (maxKeyCode+bitsPerByte)/bitsPerByte)

	request := eviocgbit(evKey, len(buf))

	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(fd),
		request,
		uintptr(unsafe.Pointer(&buf[0])),
	)

	if errno != 0 {
		return false, errno
	}

	required := []int{
		keyA,
		keyZ,
		keyEnter,
		keySpace,
	}

	for _, code := range required {
		if !bitIsSet(buf, code) {
			return false, nil
		}
	}

	return true, nil
}

func bitIsSet(bits []byte, bit int) bool {
	byteIndex := bit / 8
	bitIndex := bit % 8

	if byteIndex >= len(bits) {
		return false
	}

	return bits[byteIndex]&(1<<bitIndex) != 0
}

func deviceName(fd int) (string, error) {
	buf := make([]byte, 256)

	request := eviocgname(len(buf))

	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(fd),
		request,
		uintptr(unsafe.Pointer(&buf[0])),
	)

	if errno != 0 {
		return "", errno
	}

	for i, b := range buf {
		if b == 0 {
			return string(buf[:i]), nil
		}
	}

	return string(buf), nil
}

const (
	iocNrBits   = 8
	iocTypeBits = 8
	iocSizeBits = 14

	iocNrShift   = 0
	iocTypeShift = iocNrShift + iocNrBits
	iocSizeShift = iocTypeShift + iocTypeBits
	iocDirShift  = iocSizeShift + iocSizeBits

	iocRead = 2
)

func ioc(dir, typ, nr, size uintptr) uintptr {
	return (dir << iocDirShift) |
		(typ << iocTypeShift) |
		(nr << iocNrShift) |
		(size << iocSizeShift)
}

func ior(typ, nr, size uintptr) uintptr {
	return ioc(iocRead, typ, nr, size)
}

func eviocgbit(ev, length int) uintptr {
	return ior(
		uintptr('E'),
		uintptr(0x20+ev),
		uintptr(length),
	)
}

func eviocgname(length int) uintptr {
	return ior(
		uintptr('E'),
		0x06,
		uintptr(length),
	)
}
