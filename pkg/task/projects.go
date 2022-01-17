package task

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
)

type entityWithProject struct {
	Project string `json:"project"`
}

func GetReadyProjects() ([]string, error) {
	out, err := exec.Command("task", "export", "ready").Output()
	if err != nil {
		return nil, fmt.Errorf("export task projects: %w", err)
	}

	var es []entityWithProject
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

// GetProjectsWithNotes lists files in your configured project files location.
// Meaning projects without tasks in are not lost.
func GetProjectsWithNotes() ([]string, error) {
	projectNotesLocation, err := getConfigSetting("project_files_location")
	if err != nil {
		return []string{}, fmt.Errorf("getConfigSetting: %w", err)
	}

	// not set in .taskrc
	if projectNotesLocation == "" {
		return []string{}, nil
	}

	var fileNames []string
	filepath.Walk(projectNotesLocation, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking project notes: %w", err)
		}
		if info.IsDir() {
			return nil
		}
		name := strings.TrimPrefix(path, projectNotesLocation)
		if name[0] == '/' {
			name = name[1:]
		}
		fileNames = append(fileNames, name)
		return nil
	})

	return fileNames, nil
}

// returns the file path of project notes
func GetProjectNotesLocation(projectName string) (string, error) {
	projectNotesLocation, err := getConfigSetting("project_files_location")
	if err != nil {
		return "", fmt.Errorf("getConfigSetting: %w", err)
	}

	return fmt.Sprintf("%s/%s", projectNotesLocation, projectName), nil
}
