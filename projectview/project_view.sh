#!/bin/bash

PROJECT_NOTES_DIR="/Users/sam/Desktop/gtd/project_notes"
PROJ=$1
tmux new-window

PROJECT_FILE="$PROJECT_NOTES_DIR/$PROJ.md"
touch $PROJECT_FILE

tmux split-window -v 
tmux send-keys -t 1 "clear && cat ${PROJECT_FILE}" Enter
tmux resize-pane -U 10

pro.is does an exact match
tmux send-keys -t 0 "while true ; clear; echo "TASKS" && t pro.is:$PROJ ; sleep 5; end" Enter

tmux split-window -v
tmux resize-pane -D 10
tmux send-keys -t 2 'clear' Enter
