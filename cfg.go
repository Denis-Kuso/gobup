package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// read cfg
// validate Cfg
// flags handled elsewhere?
// cfg composed of pipelines/executers
type Cfg map[string]Pipeline

type Pipeline struct {
	Run      bool              `yaml:"run"`
	FailFast bool              `yaml:"fail_fast"`
	Steps    []map[string]hook `yaml:"cmds"`
}

type hook struct {
	Name string   `yaml:"cmdName"`
	Args []string `yaml:"args"`
}

// TODO Make pipeline

// read config file
func LoadCfg(in io.Reader) (Cfg, error) {
	var buf []byte
	buf, err := io.ReadAll(in)
	if err != nil {
		fmt.Println("could not read from %v", in)
		// TODO perhaps your own custom error type
		return Cfg{}, err
	}
	c := Cfg{}
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		fmt.Printf("config read err: %v\n", err)
		// TODO perhaps your own custom error type
		return c, err
	}
	return c, nil
}

// flags, options decide which action/pipeline will be ran
func getPipelines(c Cfg) map[string]Pipeline {
	return c
}

// TODO where does the map m comes from?
// TODO should project be passed here?
func makePipe(m []map[string]hook, project string) []executer {
	pipe := make([]executer, len(m))
	for i, step := range m {
		for name, options := range step {
			stepic := NewStep(name, options.Name, options.Args, fmt.Sprintf("%q: %s.", name, "SUCCESS"), project)
			pipe[i] = stepic
		}
	}
	return pipe
}

// print to out - mostly for debugging purposes
func printCfg(cfg *Cfg, out io.Writer) {
	for pipename, content := range cfg.Content {
		fmt.Fprintf(out, "######################\n")
		fmt.Fprintf(out, "found pipename: %q\n- properties:\n - run: %v\n - fail fast: %v\n", pipename, content.Run, content.FailFast)
		fmt.Fprintf(out, " %q has following steps:\n", pipename)
		for name, vals := range content.Steps {
			fmt.Fprintf(out, "step: %d\n", name)
			for stepName, hook := range vals {
				fmt.Fprintf(out, "  - step name: %s\n  - exe: %q\n  - args: %v\n", stepName, hook.Name, hook.Args)
			}
		}
	}
}
