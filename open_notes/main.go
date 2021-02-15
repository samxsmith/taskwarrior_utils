package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	id := os.Args[1]
	err := openTaskNotes(id)

	if err == nil {
		// done
		return
	}

	if !strings.Contains(err.Error(), "No compatible annotation found") {
		println("unable to open task: ", err.Error())
		return
	}

	err = addNotesAnnotation(id)
	if err != nil {
		println("unable to add notes annotation: ", err.Error())
		return
	}

	err = openTaskNotes(id)
	if err != nil {
		println("unable to open task: ", err.Error())
		return
	}

}

func openTaskNotes(id string) error {
	cmd := exec.Command("taskopen", id, `\\notes`)

	var cmdErr bytes.Buffer
	cmd.Stderr = &cmdErr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(cmdErr.String())
	}
	return nil
}

func addNotesAnnotation(id string) error {
	return exec.Command("task", id, "annotate", "notes: Notes").Run()
}
