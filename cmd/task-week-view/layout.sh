#!/bin/bash

tmux split-window -h 'while true ; clear; echo "URGENT" && t soon ; sleep 5; end'
tmux split-window -v 'while true ; clear; echo "PROJECTS" && t summary ; sleep 5; end'

# move back to original pane
tmux select-pane -L

tmux send-keys -t 0 'while true ; clear; echo "IMPORTANT" && t sprint ; sleep 5; end' Enter

# 
tmux split-window -v 
tmux send-keys -t 1 'clear' Enter

tmux rename-window 'gtd'

