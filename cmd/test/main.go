package main

import (
	"fmt"

	"github.com/samxsmith/taskwarriorutils/pkg/task"
)

func main() {
	fmt.Println(task.GetProjectsWithNotes())
}
