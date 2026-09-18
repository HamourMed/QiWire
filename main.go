package main

import (
	"log"

	tea "charm.land/bubbletea/v2"

	"qiwire/internal/linuxinput"
	"qiwire/internal/terminal"
	"qiwire/internal/xkb"
)

func main() {
	manager, err := linuxinput.NewKeyboardManager()
	if err != nil {
		log.Fatal(err)
	}
	defer manager.Close()

	translator, err := xkb.NewXKBTranslator("fr")
	if err != nil {
		log.Fatal(err)
	}
	defer translator.Close()

	model := terminal.New()
	program := tea.NewProgram(model)

	manager.Start()

	go func() {
		for event := range manager.Events() {
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
