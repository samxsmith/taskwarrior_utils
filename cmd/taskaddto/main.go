package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/samxsmith/taskwarriorutils/pkg/task"
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
