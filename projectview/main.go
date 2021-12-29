package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

const projectNotesDir = "/Users/sam/Desktop/gtd/project_notes"

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		return
	}
	return
}

func run() error {
	projectName, err := getProject()
	if err != nil {
		return fmt.Errorf("getProject: %w", err)
	}

	err = launch(projectName)
	if err != nil {
		return fmt.Errorf("launch: %w", err)
	}

	return nil
}

func launch(projectName string) error {
	// reset window
	if err := tmuxResetWindow(); err != nil {
		return fmt.Errorf("tmuxResetWindow: %w", err)
	}

	// setup window
	projectNotesFile := fmt.Sprintf("%s/%s", projectNotesDir, projectName)
	if err := touchFile(projectNotesFile); err != nil {
		return fmt.Errorf("touchFile) %w", err)
	}

	if err := setupWindow(projectName, projectNotesFile); err != nil {
		return fmt.Errorf("setupWindow) %w", err)
	}

	return nil
}

// clears window down to single pane
func tmuxResetWindow() error {
	// whatever the current pane is, we preserve and reuse
	// as we know that's got nothing else running in it
	// kill the rest

	// kill rest of panes
	for {
		if paneCount, err := tmuxCountPanes(); err != nil {
			return fmt.Errorf("ERROR) countPanes: %w", err)
		} else if paneCount == 1 {
			break
		}

		// each time a pane is killed, pane numbers change
		// so have to re-check current pane number each time

		currentPane, err := tmuxGetCurrentPane()
		if err != nil {
			return fmt.Errorf("tmuxGetCurrentPane: %w", err)
		}

		// kill 0 or 1 as those are the only panes we know still exist
		paneToKill := 0
		if currentPane == 0 {
			paneToKill = 1
		}
		if err := tmuxKillPane(paneToKill); err != nil {
			return fmt.Errorf("tmuxKillPane: %w", err)
		}

	}

	return nil
}

func setupWindow(projectName, projectNotesFile string) error {

	/**
	NAME PANE
	*/

	if err := tmuxNameWindow("proj:" + projectName); err != nil {
		return fmt.Errorf("tmuxNameWindow: %w", err)
	}

	/***
	PROJECT NOTES
	*/
	if err := tmuxSplitVertical(); err != nil {
		return fmt.Errorf("tmuxSplitHorizontal) %w", err)
	}

	if err := tmuxPaneCommand(1, fmt.Sprintf("clear && cat %s", projectNotesFile)); err != nil {
		return fmt.Errorf("tmuxPaneCommand) %w", err)
	}

	if err := tmuxResizeUp(10); err != nil {
		return fmt.Errorf("tmuxResizeUp) %w", err)
	}

	/***
	TASKS
	*/

	if err := tmuxPaneCommand(0, fmt.Sprintf(`while true ; clear; echo "TASKS: %s" && t pro.is:%s ; sleep 5; end`, projectName, projectName)); err != nil {
		return fmt.Errorf("tmuxPaneCommand) %w", err)
	}

	/**
	EMPTY SHELL
	*/

	if err := tmuxSplitVertical(); err != nil {
		fmt.Println("SPLIT EMPTY SHELL")
		return fmt.Errorf("tmuxSplitHorizontal) %w", err)
	}

	if err := tmuxPaneCommand(2, `clear && echo SHELL`); err != nil {
		return fmt.Errorf("tmuxPaneCommand) %w", err)
	}

	/**
	BURNDOWN
	*/

	if err := tmuxSelectPane(1); err != nil {
		return fmt.Errorf("tmuxSelectPane) %w", err)
	}

	if err := tmuxSplitHorizontal(); err != nil {
		return fmt.Errorf("tmuxSplitHorizontal) %w", err)
	}

	if err := tmuxPaneCommand(2, fmt.Sprintf(`t burndown.weekly pro:.is:%s`, projectName)); err != nil {
		return fmt.Errorf("tmuxPaneCommand) %w", err)
	}

	if err := tmuxResizeDown(10); err != nil {
		return fmt.Errorf("tmuxResizeDown) %w", err)
	}

	/**
	PROJECT SELECTOR
	*/

	if err := tmuxSelectPane(0); err != nil {
		return fmt.Errorf("tmuxSelectPane) %w", err)
	}

	if err := tmuxSplitHorizontal(); err != nil {
		return fmt.Errorf("tmuxSplitVertical) %w", err)
	}

	if err := tmuxPaneCommand(1, "task viewproject"); err != nil {
		return fmt.Errorf("tmuxPaneCommand: %w", err)
	}

	/*
		RETURN TO SHELL PANE
	*/
	if err := tmuxSelectPane(4); err != nil {
		return fmt.Errorf("tmuxSelectPane) %w", err)
	}

	return nil
}

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

func tmuxNameWindow(name string) error {
	return exec.Command("tmux", "rename-window", name).Run()
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

type EntityWithProject struct {
	Project string `json:"project"`
}

func getProject() (string, error) {
	projects, err := getAllProjects()
	if err != nil {
		return "", fmt.Errorf("ERR: getAllProjects) %w", err)
	}

	pro, err := projectSelect(projects)
	if err != nil {
		return "", fmt.Errorf("ERR: projectSelect) %w", err)
	}

	return pro, nil
}

func getAllProjects() ([]string, error) {
	out, err := exec.Command("task", "export", "ready").Output()
	if err != nil {
		return nil, fmt.Errorf("export task projects: %w", err)
	}

	var es []EntityWithProject
	err = json.Unmarshal(out, &es)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	projectsMap := map[string]bool{}
	for _, e := range es {
		if e.Project != "" {
			projectsMap[e.Project] = true
		}
	}

	projects := make([]string, len(projectsMap))

	i := 0
	for p := range projectsMap {
		projects[i] = p
		i++
	}

	// alphabetise
	sort.Strings(projects)

	return projects, nil
}

func projectSelect(projects []string) (string, error) {
	p := promptui.Select{
		Label:             "Select Project",
		Items:             projects,
		StartInSearchMode: true,
		Searcher: func(input string, index int) bool {
			if input == "" {
				return true
			}
			item := projects[index]
			return strings.Contains(item, input)
		},
	}

	_, result, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("p.Run: %w", err)
	}

	return result, nil
}

func touchFile(file string) error {
	return exec.Command("touch", file).Run()
}
