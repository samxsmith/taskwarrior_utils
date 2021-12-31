# Task AddTo

This tool lets you quick-add multiple notes to an existing project without needing to type out the project's full name.

## Dependencies
- `taskwarrior`: https://github.com/GothenburgBitFactory/taskwarrior
- `go(lang)`

## Install
```sh
go install github.com/samxsmith/taskwarrior_utils/taskaddto@latest
```

### Verify
```sh
which taskaddto
```

## Use with Taskwarrior
Add the following to your `.taskrc`:

```sh
alias.addto=execute bash -c "taskaddto"
```

This will give you a search box. Search for the project you want to add to, and press enter.

You'll be given a prompt. Start writing new taskwarrior tasks, and they'll be added to the project you've selected.
