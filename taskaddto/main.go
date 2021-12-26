package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

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

	for {
		fmt.Println("Add task to project: ", pro)
		p := promptui.Prompt{
			Label: "> ",
		}
		taskString, err := p.Run()
		if err != nil {
			if err.Error() == "^C" {
				fmt.Println("Done")
				return
			}

			fmt.Printf("ERR: %v %T \n", err, err)
			return
		}

		err = createTask(pro, taskString)
		if err != nil {
			fmt.Println("ERR createTask: ", err)
			return
		}
	}
}

type EntityWithProject struct {
	Project string `json:"project"`
}

func getAllProjects() ([]string, error) {
	out, err := exec.Command("task", "export", "ready").Output()
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

func createTask(project, taskString string) error {
	proTag := fmt.Sprintf("pro:%s", project)
	args := []string{"add", proTag}

	taskCmds := strings.Split(taskString, " ")
	args = append(args, taskCmds...)

	out, err := exec.Command("task", args...).Output()
	if err != nil {
		return fmt.Errorf("exec.Command: %w", err)
	}

	fmt.Println(string(out))
	return nil
}
