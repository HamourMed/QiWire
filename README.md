# QiWire

QiWire is an experimental Linux keylogger / keyboard input watcher written in Go.

The project is primarily a systems-programming and software-engineering exercise focused on understanding the Linux input stack, keyboard event processing, XKB translation, concurrency, and terminal user interfaces.

## Current prototype

The current version:

* scans `/dev/input/event*` and detects keyboard-capable input devices
* distinguishes typing-capable keyboard interfaces from media-only input interfaces using evdev capability queries
* opens multiple detected keyboards simultaneously
* reads raw Linux evdev events from each keyboard
* merges keyboard events through Go goroutines and channels
* converts Linux-specific input events into internal key events
* uses `libxkbcommon` to translate keycodes according to a configured keyboard layout
* handles modifiers such as Shift and Caps Lock through XKB state
* supports key repeat
* reconstructs printable UTF-8 text
* handles Enter and Backspace
* displays captured input in a terminal UI

## Architecture

```text
/dev/input/event*
        ↓
Keyboard Detector
        ↓
Typing-capable keyboards
        ↓
LinuxInputProvider × N
        ↓
goroutines
        ↓
merged KeyEvent channel
        ↓
XKBTranslator
        ↓
UTF-8 text
        ↓
TerminalDisplay
```

The current prototype targets Linux x86-64 and uses a hard-coded XKB keyboard layout.

## Build

Requirements:

* Go
* `libxkbcommon`
* `pkg-config`

Build with:

```bash
go build -o qiwire .
```

Run with sufficient permissions to access Linux input devices:

```bash
sudo ./qiwire
```

## Status

QiWire is experimental and under active development.

The current implementation focuses on validating the core input pipeline before adding features such as:

* active Wayland keyboard-layout detection
* keyboard hotplug and disconnect handling
* cleaner goroutine cancellation and shutdown
* more robust physical-device identification
* Unicode-safe editing behavior
* additional platform backends for Windows and macOS

## Disclaimer

QiWire is intended for educational, experimental, and authorized security research use only. Keyboard input can contain highly sensitive information, so the tool should only be used on systems and input devices you are authorized to monitor.
