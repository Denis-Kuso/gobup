package main

import (
	"fmt"
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
	Name     string // ideally this would be the name of the pipeline
	Steps    []map[string]hook `yaml:"cmds"`
	Run      bool   `yaml:"run"`
	FailFast bool   `yaml:"fail_fast"`
}

type hook struct {
	Name string `yaml:"cmdName"`
	Args []string `yaml:"args"`
}


func loadCfg() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("no config file found, probably")
		os.Exit(23)
	}
	draft0 := Cfg{}
	err = yaml.Unmarshal(file, &draft0.Content)
	if err != nil {
		fmt.Printf("config read err: %v\n", err)
	}
	for pipename, content := range draft0.Content {
		fmt.Printf("######################\n")
		fmt.Printf("found pipename: %q\n- properties:\n - run: %v\n - fail fast: %v\n", pipename, content.Run, content.FailFast)
		fmt.Printf("%s has following steps:\n", pipename)
		for name, vals := range content.Steps {
			fmt.Printf("step: %d\n", name)
			for stepName, hook := range vals {
				fmt.Printf("  - step name: %s\n  - exe: %q\n  - args: %v\n", stepName, hook.Name, hook.Args)
			}
		}
	}
}
