# QiWire

QiWire is an experimental Linux keylogger / keyboard input watcher written in Go.

The project is primarily a systems-programming and software-engineering exercise focused on understanding the Linux input stack, keyboard event processing, XKB translation, and terminal user interfaces.

## Current prototype

The current version:

- reads keyboard events directly from `/dev/input/eventX` using the Linux evdev interface
- converts raw Linux input events into internal key events
- uses `libxkbcommon` to translate keycodes according to a configured keyboard layout
- handles modifiers such as Shift and Caps Lock through XKB state
- supports key repeat
- reconstructs printable UTF-8 text
- handles Enter and Backspace
- displays captured input in a terminal UI

## Architecture

```text
/dev/input/eventX
        ↓
LinuxInputProvider
        ↓
KeyEvent
        ↓
XKBTranslator
        ↓
UTF-8 text
        ↓
TerminalDisplay
```

The project currently targets Linux x86-64 and uses a manually selected input device and keyboard layout.

## Build

Requirements:

- Go
- `libxkbcommon`
- `pkg-config`

Build with:

```bash
go build -o qiwire .
```

Run with sufficient permissions to read the selected `/dev/input/eventX` device:

```bash
sudo ./qiwire
```

## Status

QiWire is experimental and under active development. The current implementation is intentionally small and focuses on validating the core input and translation pipeline before adding features such as automatic keyboard discovery, hotplug support, active Wayland layout detection, and additional platform backends.

## Disclaimer

QiWire is intended for educational, experimental, and authorized security research use only. Keyboard input can contain highly sensitive information, so the tool should only be used on systems and input devices you are authorized to monitor.