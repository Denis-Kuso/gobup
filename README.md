## Workflow 1.0 - command description

### "set and forget" usage: git users

```bash
./gobup init --git .
```

Creates cfg file and creates pre-commit, pre-push files (must be git repo).
Upon every attempt to commit, commands specified in the cfg file associated
with the action `pre-commit` will be ran. By default any errors with the
specified commands will be treated as warnings (error in step 1 does not
prevent running of step 2). This can be changed with the `-e` flag or modifying
the option in the cfg file.

Upon every attempt to push, the set of commands associated with `pre-push` will
run. This action treats all warnings as errors by default (if step 1 fails,
step 2 will not run).

### Example/default cfg

```yaml
# pipeline name
pre-commit:
  run: true
  # if possible, all commands run as warnings
  fail_fast: false
  # sequence of commands to run in this pipeline
  cmds:
    - build:
        cmdName: go build
        # args are ordered
        args:
          - "-o"
          - "binaryName"
    - test:
        cmdName: go test
        args:
          - "."
          - "-v"
    - format:
        cmdName: gofmt
        args: ["-l", "."]
pre-push:
  run: false
  # failed command will prevent execution of next command
  fail_fast: true
  cmds:
    # commands run in specified order (build, test, format, push)
    - build:
        cmdName: go build
        args:
          - "."
    - test:
        cmdName: go test
        args:
          - "."
          - "-v"
    - format:
        cmdName: gofmt
        # as per YAML specification, can also specify order using brackets
        args: ["-l", "."]
    - push:
        cmdName: "git push"
        args:
          - "origin"
          - "main"
        # will fail if command takes longer than timeout seconds
        timeout: 5
```

### A bit more manual usage - no git installed

```bash
./gobup init .
```
Creates cfg file local to project, which may or may not be a git repo.

```bash
./gobup run
```
runs all actions specified in the cfg file in current directory, which have
the propery `run: true`. Invoking this command is not neccessary if
initialised with `--git` option specified.

```bash
./gobup run --dry-run
```
Prints to stdout all the actions - and hooks that would be ran.

```bash
./gobup run -e
```

Runs all actions with warnings treated as errors (if possible).

```bash
./gobup run -p <pipeline>
```

Runs only the commands asociated with the `<action_name>` regardles if
`run: false` specified in cfg file (acts like a singular overwrite).
