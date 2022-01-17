# Taskwarrior Utils

Taskwarrior is awesome.

The best thing about it is that you can build on top of it, to add all the functionality you want.

Here you'll find all the tools I've added to make it work the way I'd like.

## Installation
### Dependencies
The `view` tools depend upon tmux. If you don't want to use them, the install will work fine.

All depend upon having taskwarrior installed with the binary name `task`.

## Tools
### Open
The popular taskwarrior `open` script is a little hard to install, and written in Perl. This is a simplified version written in Go. It is not as fully featured as the Perl version, but it does the job.

[Do check out the real taskopen repo here.](https://github.com/jschlatow/taskopen)

Install our simplified version by running:
`go install github.com/samxsmith/taskwarrior_utils/taskopen`

### Open Notes
A simple tool to give you notes per task in one command. e.g.
```
task notes 12
```
[Find out more](cmd/open_notes/)

### Add to
Adding a task to an existing project can be a pain. Sometimes it's nice to have longer, explanatory project names.
Especially when working with sub-projects.

Instead of having to copy and paste the project name this tool gives you a select capability, to search for a project
and then add multiple tasks once selected.

[Find out more](cmd/taskaddto/)

### Project View
Using `tmux`, this sets up a fully project view, including easy dropdown switcher, your burndown and tasks.

#### Dependencies
- `tmux`

### Week View
Using `tmux`, lays out an urgent and important panes.

The important pane looks at the `planned` field. If you don't have this configured it won't work.

The urgent pane looks for tasks due in the next 7 days.

# More to come
I'll slowly be pushing my utils here as I tidy them up.
