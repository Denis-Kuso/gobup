# gobup

## What

A configuration-based local pipeline builder (CI-like). Edit a config file,
and have different commands execute based on your needs/workflow.

## Why

Previously used a combination of make files, one-of python/bash scripts to run
different commands (testing, formating, etc).
Sometimes made/edited hook files. All of which somewhat cumbersome for maintenance.
Wanted something that:

- is CLI-based
- can run independent of IDE/plugins.
- can be used easily with git hook files in git

Also, whilst developing, you might not need to adhere to
the same rules when commiting/pushing to a repo, hence can have a different set
of commands executed or in different order (same goes for different projects).

## Installation

Using go:

`go install github.com/Denis-Kuso/gobup`

or from [releases](github.com/Denis-Kuso/gobup/releases/) page.

## Usage

```bash
gobup init <project_dir>
```
Creates template cfg file in the provided directory named `.gobup.yaml`. Best to provide
the root of your project.
See example below:

<details markdown=1><summary markdown="span">Example cfg</summary>

```yaml
# pipeline name
pre-commit:
  run: false
  # sequence (order matters) of commands to run in this pipeline
  cmds:
    - build:
        cmdName: go
        # args to cmdName (ordered)
        args:
          - "build"
          - "-o"
          - "binaryName"
    - test:
        cmdName: go
        args:
          - "test"
          - "."
          - "-v"
    - format:
        cmdName: gofmt
        args: ["-l", "."]
        # output to stdout interpreted/treated as an error
        stdoutAsErr: true
#  different pipeline perhaps ran whilst prototyping/developing
dev:
  run: true
  cmds:
    - lint:
        cmdName: revive
        args:
          - "-formatter"
          - "friendly"
        # stop execution if it takes longer than timeout seconds
        timeout: 15
    - format:
        cmdName: gofmt
        args: ["-l", "."]
        stdoutAsErr: true
```
</details>

```bash
gobup run
```

Runs all actions specified in the cfg file in current directory, which have
the propery `run: true`.

```bash
gobup run -p <pipeline>
```

Only run the commands associated with `<pipeline>`, regardles if
`run: false` for that pipeline. Ignores other pipelines.

## build


<details markdown=1><summary markdown="span">TODO</summary>

- [ ] add `dry-run` flag
- [ ] add `ignore-warnings` flag
- [ ] make prettier output format
- [ ] add git hook files compatibility
</details>
