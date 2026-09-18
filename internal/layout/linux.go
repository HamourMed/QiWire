package layout

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	defaultKeyboardPath = "/etc/default/keyboard"
	vconsolePath        = "/etc/vconsole.conf"
)

type Config struct {
	Model   string
	Layout  string
	Variant string
	Options string
}

func Detect() (Config, error) {
	// Prefer the XKB configuration because it maps directly
	// to what libxkbcommon expects.
	config, err := detectDefaultKeyboard()
	if err == nil {
		return config, nil
	}

	// Fallback to the virtual-console configuration.
	config, vconsoleErr := detectVConsole()
	if vconsoleErr == nil {
		return config, nil
	}

	return Config{}, fmt.Errorf(
		"failed to detect keyboard layout: %v; fallback failed: %v",
		err,
		vconsoleErr,
	)
}

func detectDefaultKeyboard() (Config, error) {
	file, err := os.Open(defaultKeyboardPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = cleanValue(value)

		switch key {
		case "XKBMODEL":
			config.Model = value

		case "XKBLAYOUT":
			config.Layout = value

		case "XKBVARIANT":
			config.Variant = value

		case "XKBOPTIONS":
			config.Options = value
		}
	}

	if err := scanner.Err(); err != nil {
		return Config{}, err
	}

	if config.Layout == "" {
		return Config{}, fmt.Errorf(
			"XKBLAYOUT not found in %s",
			defaultKeyboardPath,
		)
	}

	return config, nil
}

func detectVConsole() (Config, error) {
	file, err := os.Open(vconsolePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = cleanValue(value)

		if key == "KEYMAP" {
			if value == "" {
				return Config{}, fmt.Errorf("KEYMAP is empty")
			}

			return Config{
				Layout: value,
			}, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return Config{}, err
	}

	return Config{}, fmt.Errorf(
		"KEYMAP not found in %s",
		vconsolePath,
	)
}

func cleanValue(value string) string {
	value = strings.TrimSpace(value)
	return strings.Trim(value, `"'`)
}
