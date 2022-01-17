package main

import (
	"fmt"

	"github.com/samxsmith/taskwarriorutils/pkg/tmux"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("ERROR) %s \n", err)
	}
}

func run() error {
	if err := tmux.SplitHorizontal(); err != nil {
		return fmt.Errorf("tmux.SplitHorizontal: %w", err)
	}
	if err := tmux.SplitVertical(); err != nil {
		return fmt.Errorf("tmux.SplitVertical: %w", err)
	}
	if err := tmux.SelectPane(0); err != nil {
		return fmt.Errorf("tmux.SelectPane: %w", err)
	}
	if err := tmux.SplitVertical(); err != nil {
		return fmt.Errorf("tmux.SplitVertical: %w", err)
	}

	if err := tmux.PaneCommand(0, `while true ; clear; echo "IMPORTANT" && t sprint ; sleep 5; end`); err != nil {
		return fmt.Errorf("tmux.PaneCommand: %w", err)
	}
	if err := tmux.PaneCommand(1, `clear`); err != nil {
		return fmt.Errorf("tmux.PaneCommand: %w", err)
	}
	if err := tmux.PaneCommand(2, `while true ; clear; echo "URGENT" && t soon ; sleep 5; end`); err != nil {
		return fmt.Errorf("tmux.PaneCommand: %w", err)
	}
	if err := tmux.PaneCommand(3, `while true ; clear; echo "PROJECTS" && t summary ; sleep 5; end`); err != nil {
		return fmt.Errorf("tmux.PaneCommand: %w", err)
	}

	if err := tmux.SetWindowName("gtd"); err != nil {
		return fmt.Errorf("tmux.SetWindowName: %w", err)
	}

	return nil
}
