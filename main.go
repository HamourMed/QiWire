package main

import (
	"log"

	tea "charm.land/bubbletea/v2"

	"qiwire/internal/linuxinput"
	"qiwire/internal/terminal"
	"qiwire/internal/xkb"
)

const inputDevice = "/dev/input/event7"

func main() {
	provider, err := linuxinput.NewLinuxInputProvider(inputDevice)
	if err != nil {
		log.Fatal(err)
	}
	defer provider.Close()

	translator, err := xkb.NewXKBTranslator("fr")
	if err != nil {
		log.Fatal(err)
	}
	defer translator.Close()

	model := terminal.New()
	program := tea.NewProgram(model)

	go func() {
		for {
			event, err := provider.ReadEvent()
			if err != nil {
				program.Quit()
				return
			}

			text := translator.Translate(event)
			if text != "" {
				program.Send(terminal.TextMsg(text))
			}
		}
	}()

	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
