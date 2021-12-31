package tmux

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func CountPanes() (int, error) {
	paneList, err := exec.Command("tmux", "list-panes").Output()
	if err != nil {
		return 0, fmt.Errorf("list-panes: %w", err)
	}

	count := strings.Count(string(paneList), "\n")
	return count, nil
}

func NewWindow() error {
	return exec.Command("tmux", "new-window").Run()
}

func SelectPane(paneNum int) error {
	return exec.Command("tmux", "select-pane", "-t", strconv.Itoa(paneNum)).Run()
}
func KillPane(paneNum int) error {
	return exec.Command("tmux", "kill-pane", "-t", strconv.Itoa(paneNum)).Run()
}
func GetCurrentPane() (int, error) {
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

func SplitHorizontal() error {
	return Split("-h")
}
func SplitVertical() error {
	return Split("-v")
}

func Split(directionFlag string) error {
	return exec.Command("tmux", "split-window", directionFlag).Run()
}

func SetWindowName(name string) error {
	return exec.Command("tmux", "rename-window", name).Run()
}
func GetWindowName() (string, error) {
	b, err := exec.Command("tmux", "display-message", "-p", `"#{window_name}"`).Output()
	if err != nil {
		return "", err
	}
	s := strings.Trim(string(b[:len(b)-1]), `""`)
	return s, nil
}

func PaneCommand(paneNum int, command string) error {
	paneNumI := strconv.Itoa(paneNum)
	return exec.Command("tmux", "send-keys", "-t", paneNumI, command, "Enter").Run()
}

func PaneKillCommand(paneNum int) error {
	paneNumI := strconv.Itoa(paneNum)
	return exec.Command("tmux", "send-keys", "-t", paneNumI, "C-c").Run()
}

func ResizeDown(amount int) error {
	return Resize(amount, "-D")
}
func ResizeUp(amount int) error {
	return Resize(amount, "-U")
}
func Resize(amount int, directionFlag string) error {
	amountS := strconv.Itoa(amount)
	return exec.Command("tmux", "resize-pane", directionFlag, amountS).Run()
}

// clears window down to single pane
func ResetWindow() error {
	// whatever the current pane is, we preserve and reuse
	// as we know that's got nothing else running in it
	// kill the rest

	// kill rest of panes
	for {
		if paneCount, err := CountPanes(); err != nil {
			return fmt.Errorf("ERROR) countPanes: %w", err)
		} else if paneCount == 1 {
			break
		}

		// each time a pane is killed, pane numbers change
		// so have to re-check current pane number each time

		currentPane, err := GetCurrentPane()
		if err != nil {
			return fmt.Errorf("GetCurrentPane: %w", err)
		}

		// kill 0 or 1 as those are the only panes we know still exist
		paneToKill := 0
		if currentPane == 0 {
			paneToKill = 1
		}
		if err := KillPane(paneToKill); err != nil {
			return fmt.Errorf("tmux.KillPane: %w", err)
		}

	}

	return nil
}
