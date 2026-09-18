package linuxinput

import (
	"fmt"
)

type KeyboardManager struct {
	providers []*LinuxInputProvider
	events    chan KeyEvent
}

func NewKeyboardManager() (*KeyboardManager, error) {
	devices, err := DetectKeyboards()
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no keyboards detected")
	}

	manager := &KeyboardManager{
		events: make(chan KeyEvent),
	}

	for _, device := range devices {
		provider, err := NewLinuxInputProvider(device.Path)
		if err != nil {
			manager.Close()
			return nil, fmt.Errorf(
				"failed to open %s (%s): %w",
				device.Name,
				device.Path,
				err,
			)
		}

		manager.providers = append(manager.providers, provider)
	}

	return manager, nil
}

func (m *KeyboardManager) Start() {
	for _, provider := range m.providers {
		go func(p *LinuxInputProvider) {
			for {
				event, err := p.ReadEvent()
				if err != nil {
					return
				}

				m.events <- event
			}
		}(provider)
	}
}

func (m *KeyboardManager) Events() <-chan KeyEvent {
	return m.events
}

func (m *KeyboardManager) Close() {
	for _, provider := range m.providers {
		_ = provider.Close()
	}
}
