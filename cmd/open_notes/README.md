# Task Open Notes

`taskopen` is an awesome tool for opening links in annotations.
You can read more about it here: [https://github.com/jschlatow/taskopen](https://github.com/jschlatow/taskopen)

The feature I use most is that if there is an annotation `notes: Notes` then it will create or open a notes file in a directory you can specify.

BUT I'm lazy, and having to add that annotation is extra commands and keystrokes.
This script will add the annotation for you if it does not exist, then open the note.

This tool is now my **most used** taskwarrior command.

## Dependencies
- `taskwarrior`: https://github.com/GothenburgBitFactory/taskwarrior
- `taskopen`: https://github.com/jschlatow/taskopen
- `go(lang)`

## Installing

```sh
# install all the utils
git clone git@github.com:samxsmith/taskwarrior_utils.git
cd open_notes
make install
```

## Using
You can use this command directly:
```
task-open-notes 12
```
This will open notes for the taskwarrior task 12.

### Using with taskwarrior
More likely you'll want to use it from taskwarrior directly.
By adding the following line to your `~/.taskrc`:
```sh
alias.notes=execute bash -l -c "q=($BASH_COMMAND); task-open-notes \\"\\\\${q[4]}\\""
```

the following command will now work:
```
task notes 12
```

Opening your task notes automatically, without any per-task setup.

This is my most-used taskwarrior command.
