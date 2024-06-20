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
type Cfg struct {
	Content map[string]Pipeline
}
type Pipeline struct {
	Name     string            // ideally this would be the name of the pipeline
	Steps    []map[string]hook `yaml:"cmds"`
	Run      bool              `yaml:"run"`
	FailFast bool              `yaml:"fail_fast"`
}

type hook struct {
	Name string   `yaml:"cmdName"`
	Args []string `yaml:"args"`
}

// TODO Make pipeline
// return a collection of pipelines?

// read config file
func loadCfg(debug bool) (Cfg, error) {
	// TODO refactor
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("no config file found, probably")
		// TODO perhaps your own custom error type
		return Cfg{}, err
	}
	draft0 := Cfg{}
	err = yaml.Unmarshal(file, &draft0.Content)
	if err != nil {
		fmt.Printf("config read err: %v\n", err)
		// TODO perhaps your own custom error type
		return draft0, err
	}
	return draft0, nil
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
