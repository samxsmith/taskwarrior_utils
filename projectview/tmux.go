package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func tmuxCountPanes() (int, error) {
	paneList, err := exec.Command("tmux", "list-panes").Output()
	if err != nil {
		return 0, fmt.Errorf("list-panes: %w", err)
	}

	count := strings.Count(string(paneList), "\n")
	return count, nil
}

func tmuxNewWindow() error {
	return exec.Command("tmux", "new-window").Run()
}

func tmuxSelectPane(paneNum int) error {
	return exec.Command("tmux", "select-pane", "-t", strconv.Itoa(paneNum)).Run()
}
func tmuxKillPane(paneNum int) error {
	return exec.Command("tmux", "kill-pane", "-t", strconv.Itoa(paneNum)).Run()
}
func tmuxGetCurrentPane() (int, error) {
	res, err := exec.Command("tmux", "display-message", "-p", `"#{pane_index}`).Output()
	if err != nil {
		return 0, err
	}

	// chop off trailing \n
	output := string(res[:len(res)-1])
	paneNum, err := strconv.Atoi(strings.Trim(output, `"\n`))
	if err != nil {
		return 0, err
	}
	return paneNum, nil
}

func tmuxSplitHorizontal() error {
	return tmuxSplit("-h")
}
func tmuxSplitVertical() error {
	return tmuxSplit("-v")
}

func tmuxSplit(directionFlag string) error {
	return exec.Command("tmux", "split-window", directionFlag).Run()
}

func tmuxSetWindowName(name string) error {
	return exec.Command("tmux", "rename-window", name).Run()
}
func tmuxGetWindowName() (string, error) {
	b, err := exec.Command("tmux", "display-message", "-p", `"#{window_name}"`).Output()
	if err != nil {
		return "", err
	}
	s := strings.Trim(string(b[:len(b)-1]), `""`)
	return s, nil
}

func tmuxPaneCommand(paneNum int, command string) error {
	paneNumI := strconv.Itoa(paneNum)
	return exec.Command("tmux", "send-keys", "-t", paneNumI, command, "Enter").Run()
}

func tmuxPaneKillCommand(paneNum int) error {
	paneNumI := strconv.Itoa(paneNum)
	return exec.Command("tmux", "send-keys", "-t", paneNumI, "C-c").Run()
}

func tmuxResizeDown(amount int) error {
	return tmuxResize(amount, "-D")
}
func tmuxResizeUp(amount int) error {
	return tmuxResize(amount, "-U")
}
func tmuxResize(amount int, directionFlag string) error {
	amountS := strconv.Itoa(amount)
	return exec.Command("tmux", "resize-pane", directionFlag, amountS).Run()
}
