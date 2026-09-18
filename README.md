# QiWire

QiWire is an experimental Linux keylogger / keyboard input watcher written in Go.

The project is primarily a systems-programming and software-engineering exercise focused on understanding the Linux input stack, keyboard event processing, XKB translation, concurrency, and terminal user interfaces.

## Current prototype

The current version:

* scans `/dev/input/event*` and detects typing-capable keyboard devices
* distinguishes normal keyboard interfaces from media-only input interfaces using evdev capability queries
* opens multiple detected keyboards simultaneously
* reads raw Linux evdev events from each keyboard
* merges keyboard events through Go goroutines and channels
* converts Linux-specific events into internal `KeyEvent` values
* detects the system XKB configuration from Linux keyboard configuration files
* loads XKB model, layout, variant, and options into `libxkbcommon`
* handles modifiers such as Shift, Caps Lock, and AltGr through XKB state
* supports key repeat
* reconstructs printable UTF-8 text
* handles Enter and Backspace
* displays captured input in a scrollable terminal UI

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

Keyboard configuration is currently discovered from the system configuration, preferring `/etc/default/keyboard` when available and falling back to `/etc/vconsole.conf`.

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

## Current limitations

QiWire currently targets Linux x86-64.

The current prototype does not yet handle:

* active Wayland layout changes at runtime
* keyboard hotplug and disconnect recovery
* graceful cancellation of all input goroutines
* robust physical-device grouping
* full cross-platform input backends
* complete console-keymap to XKB-layout conversion for every Linux configuration

## Roadmap

Planned improvements include:

* active Wayland keyboard-layout detection
* keyboard hotplug support
* cleaner shutdown and error propagation
* better physical keyboard identification
* improved Unicode editing behavior
* Windows and macOS input backends

## Disclaimer

QiWire is intended for educational, experimental, and authorized security research use only. Keyboard input can contain highly sensitive information, so the tool should only be used on systems and input devices you are authorized to monitor.
