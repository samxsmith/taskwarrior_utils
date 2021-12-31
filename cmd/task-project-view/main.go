package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/samxsmith/taskwarriorutils/pkg/task"
	"github.com/samxsmith/taskwarriorutils/pkg/tmux"
)

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
	if err := tmux.ResetWindow(); err != nil {
		return fmt.Errorf("tmux.ResetWindow: %w", err)
	}

	// ensure project notes exist
	projectNotesFile, err := task.GetProjectNotesLocation(projectName)
	if err != nil {
		return fmt.Errorf("task.GetProjectNotesLocation) %w", err)
	}

	if err := touchFile(projectNotesFile); err != nil {
		return fmt.Errorf("touchFile) %w", err)
	}

	// setup window
	if err := setupWindow(projectName, projectNotesFile); err != nil {
		return fmt.Errorf("setupWindow) %w", err)
	}

	return nil
}

func setupWindow(projectName, projectNotesFile string) error {

	/**
	NAME PANE
	*/

	if err := tmux.SetWindowName("proj:" + projectName); err != nil {
		return fmt.Errorf("tmux.SetWindowName: %w", err)
	}

	/***
	PROJECT NOTES
	*/
	if err := tmux.SplitVertical(); err != nil {
		return fmt.Errorf("tmux.SplitHorizontal) %w", err)
	}

	if err := tmux.PaneCommand(1, fmt.Sprintf("clear && cat %s", projectNotesFile)); err != nil {
		return fmt.Errorf("tmux.PaneCommand) %w", err)
	}

	if err := tmux.ResizeUp(10); err != nil {
		return fmt.Errorf("tmux.ResizeUp) %w", err)
	}

	/***
	TASKS
	*/

	if err := tmux.PaneCommand(0, fmt.Sprintf(`while true ; clear; echo "TASKS: %s" && t pro.is:%s ; sleep 5; end`, projectName, projectName)); err != nil {
		return fmt.Errorf("tmux.PaneCommand) %w", err)
	}

	/**
	EMPTY SHELL
	*/

	if err := tmux.SplitVertical(); err != nil {
		return fmt.Errorf("tmux.SplitHorizontal) %w", err)
	}

	if err := tmux.PaneCommand(2, `clear && echo SHELL`); err != nil {
		return fmt.Errorf("tmux.PaneCommand) %w", err)
	}

	/**
	PROJECT SELECTOR
	*/

	if err := tmux.SplitHorizontal(); err != nil {
		return fmt.Errorf("tmux.SplitVertical) %w", err)
	}

	if err := tmux.PaneCommand(3, "clear && task viewproject"); err != nil {
		return fmt.Errorf("tmux.PaneCommand: %w", err)
	}

	/**
	BURNDOWN
	*/

	if err := tmux.SelectPane(1); err != nil {
		return fmt.Errorf("tmux.SelectPane) %w", err)
	}

	if err := tmux.SplitHorizontal(); err != nil {
		return fmt.Errorf("tmux.SplitHorizontal) %w", err)
	}

	if err := tmux.PaneCommand(2, fmt.Sprintf(`t burndown.weekly pro.is:%s`, projectName)); err != nil {
		return fmt.Errorf("tmux.PaneCommand) %w", err)
	}

	/*
		RETURN TO SHELL PANE
	*/
	if err := tmux.SelectPane(3); err != nil {
		return fmt.Errorf("tmux.SelectPane) %w", err)
	}

	return nil
}

type EntityWithProject struct {
	Project string `json:"project"`
}

func getProject() (string, error) {
	projects, err := getAllProjects()
	if err != nil {
		return "", fmt.Errorf("getAllProjects) %w", err)
	}

	currentProjectName, err := getCurrentProjectName()
	if err != nil {
		return "", fmt.Errorf("getCurrentProjectName) %w", err)
	}

	pro, err := projectSelect(currentProjectName, projects)
	if err != nil {
		return "", fmt.Errorf("projectSelect) %w", err)
	}

	return pro, nil
}

func getCurrentProjectName() (string, error) {
	windowName, err := tmux.GetWindowName()
	if err != nil {
		return "", fmt.Errorf("tmux.GetWindowName) %w", err)
	}

	var currentProject string
	if strings.HasPrefix(windowName, "proj:") {
		currentProject = windowName[5:]
	}

	return currentProject, nil
}

func getAllProjects() ([]string, error) {
	readyProjects, err := task.GetReadyProjects()
	if err != nil {
		return []string{}, fmt.Errorf("task.GetReadyProjects: %w", err)
	}

	// check these too, in case projects don't have any tasks, and aren't forgotten
	projectsWithNotes, err := task.GetProjectsWithNotes()
	if err != nil {
		return []string{}, fmt.Errorf("task.GetProjectsWithNotes: %w", err)
	}

	allProjects := uniqueStringSlice(readyProjects, projectsWithNotes)
	// alphabetise
	sort.Strings(allProjects)
	return allProjects, nil
}

func projectSelect(currentProject string, projects []string) (string, error) {

	var currentProjectIndex int
	if currentProject != "" {
		for i, proj := range projects {
			if proj == currentProject {
				currentProjectIndex = i
				break
			}
		}
	}

	startInSearchMode := true

	// selecting current project doesn't work in search mode
	// search mode can be triggered with `/`
	if currentProjectIndex > 0 {
		startInSearchMode = false
	}

	p := promptui.Select{
		Label:             "Select Project",
		Items:             projects,
		StartInSearchMode: startInSearchMode,
		Searcher: func(input string, index int) bool {
			if input == "" {
				return true
			}
			item := projects[index]
			return strings.Contains(item, input)
		},
		Size: 15,

		// start at current project
		CursorPos: currentProjectIndex,
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

func uniqueStringSlice(a []string, rest ...[]string) []string {
	m := map[string]bool{}

	for _, v := range a {
		m[v] = true
	}
	for _, s := range rest {
		for _, v := range s {
			m[v] = true
		}
	}

	result := make([]string, len(m))

	var i int
	for k := range m {
		result[i] = k
		i++
	}

	return result
}
