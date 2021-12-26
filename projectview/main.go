package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

const projectNotesDir = "/Users/sam/Desktop/gtd/project_notes"

func main() {
	projects, err := getAllProjects()
	if err != nil {
		fmt.Println("ERR: getAllProjects) ", err)
		return
	}

	pro, err := projectSelect(projects)
	if err != nil {
		fmt.Println("ERR: projectSelect) ", err)
		return
	}

	proNotesFiles := fmt.Sprintf("%s/%s", projectNotesDir, pro)
	if err = touchFile(proNotesFiles); err != nil {
		fmt.Printf("ERROR: touchFile) %s \n", err)
		return
	}

	err = tmuxNewWindow()
	if err != nil {
		fmt.Printf("ERROR: tmuxNewWindow) %s \n", err)
		return
	}
	err = tmuxSplitVertical()
	if err != nil {
		fmt.Printf("ERROR: tmuxSplitVertical) %s \n", err)
		return
	}
	err = tmuxPaneCommand(1, fmt.Sprintf("clear && cat %s", proNotesFiles))
	if err != nil {
		fmt.Printf("ERROR: tmuxPaneCommand) %s \n", err)
		return
	}
	err = tmuxResizeUp(10)
	if err != nil {
		fmt.Printf("ERROR: tmuxResizeUp) %s \n", err)
		return
	}

	err = tmuxPaneCommand(0, fmt.Sprintf(`while true ; clear; echo "TASKS" && t pro.is:%s ; sleep 5; end`, pro))
	if err != nil {
		fmt.Printf("ERROR: tmuxPaneCommand) %s \n", err)
		return
	}

	err = tmuxSplitVertical()
	if err != nil {
		fmt.Printf("ERROR: tmuxSplitVertical) %s \n", err)
		return
	}

	err = tmuxPaneCommand(2, fmt.Sprintf(`t burndown.weekly pro:.is:%s`, pro))
	if err != nil {
		fmt.Printf("ERROR: tmuxPaneCommand) %s \n", err)
		return
	}

	err = tmuxSplitHorizontal()
	if err != nil {
		fmt.Printf("ERROR: tmuxSplitVertical) %s \n", err)
		return
	}
	// err = tmuxResizeDown(10)
	// if err != nil {
	// 	fmt.Printf("ERROR: tmuxResizeDown) %s \n", err)
	// 	return
	// }
	err = tmuxPaneCommand(3, `clear`)
	if err != nil {
		fmt.Printf("ERROR: tmuxPaneCommand) %s \n", err)
		return
	}
}

func tmuxNewWindow() error {
	return exec.Command("tmux", "new-window").Run()
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

func tmuxPaneCommand(paneNum int, command string) error {
	paneNumI := strconv.Itoa(paneNum)
	return exec.Command("tmux", "send-keys", "-t", paneNumI, command, "Enter").Run()
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

func getAllProjects() ([]string, error) {
	out, err := exec.Command("task", "export", "-COMPLETED", "-DELETED").Output()
	if err != nil {
		return nil, fmt.Errorf("exec task projects: %w", err)
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
