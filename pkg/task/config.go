package task

import (
	"os/exec"
	"strings"
)

func getConfigSetting(settingName string) (string, error) {
	b, err := exec.Command("task", "show", settingName).Output()

	output := strings.Trim(string(b), "\n ")

	lines := strings.Split(output, "\n")
	if len(lines) < 3 {
		return "", nil
	}
	// first line is header, second line is horizontal divider, third line is result
	configLine := lines[2]

	// split and ignore empty
	configParts := strings.FieldsFunc(configLine, func(r rune) bool {
		return r == ' '
	})

	if len(configParts) != 2 {
		return "", nil
	}

	// first is config name, second is config value
	return configParts[1], err
}
