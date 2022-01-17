package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var taskID int
	for _, a := range os.Args[1:] {
		if t, err := strconv.Atoi(a); err == nil {
			taskID = t
			break
		}
	}

	if taskID == 0 {
		fmt.Println("no taskID was passed")
		return
	}

	err := taskOpen(taskID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func taskOpen(taskID int) error {
	annotationCount, err := countAnnotations(taskID)
	if err != nil {
		return fmt.Errorf("could not count annotations: %v", err)
	}

	if annotationCount == 0 {
		return fmt.Errorf("no annotations to open")
	}

	if annotationCount == 1 {
		return openAnnotation(taskID, 1)
	}

	annotationIndex, err := chooseAnnotation(taskID, annotationCount)
	if err != nil {
		return fmt.Errorf("chooseAnnotation: %w", err)
	}

	return openAnnotation(taskID, annotationIndex)
}

func countAnnotations(taskID int) (int, error) {
	str, err := runExec("task", "_get", fmt.Sprintf("%d.annotations.count", taskID))
	if err != nil {
		return 0, err
	}

	c, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi: %v", err)
	}

	return c, nil
}

func openAnnotation(taskID, annotationIndex int) error {
	annotationVal, err := getAnnotationValue(taskID, annotationIndex)
	if err != nil {
		return fmt.Errorf("couldn't get annotation value: %w", err)
	}

	fmt.Printf("Opening %s... \n", annotationVal)
	_, err = runExec("open", annotationVal)
	if err != nil {
		return fmt.Errorf("couldn't open annotation: %w", err)
	}
	return nil
}

func getAnnotationValue(taskID, annotationIndex int) (string, error) {
	annotationVal, err := runExec("task", "_get", fmt.Sprintf("%d.annotations.%d.description", taskID, annotationIndex))
	if err != nil {
		return "", fmt.Errorf("couldn't get annotation value: %w", err)
	}
	return annotationVal, nil
}

func chooseAnnotation(taskID, annotationCount int) (int, error) {
	fmt.Println("Select annotation by index: ")
	for i := 0; i < annotationCount; i++ {
		// tw indexes from 1
		annotationIndex := i + 1
		annotationVal, err := getAnnotationValue(taskID, annotationIndex)
		if err != nil {
			return 0, fmt.Errorf("couldn't get annotation value: %w", err)
		}
		fmt.Printf("[%d] %s \n", annotationIndex, annotationVal)
	}

	p := promptui.Prompt{
		Label: "choose an index: ",
		Validate: func(s string) error {
			_, err := strconv.Atoi(s)
			return err
		},
	}
	str, err := p.Run()
	if err != nil {
		return 0, fmt.Errorf("error choosing annotation: %w", err)
	}
	annotationIndex, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("index passed is not an integer: %w", err)
	}

	return annotationIndex, nil
}

func runExec(command string, args ...string) (string, error) {
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", fmt.Errorf("exec.Command: %w", err)
	}

	str := string(output)

	// remove spaces and empty lines from start and end of output
	str = strings.Trim(str, " \n")

	return str, nil
}
