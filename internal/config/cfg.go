package config

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

var ErrConfig = errors.New("configuration failure")

// read cfg
// validate Cfg
// flags handled elsewhere?
// cfg composed of pipelines/executers
type Cfg map[string]Pipeline
type Action map[string]hook

type Pipeline struct {
	Run   bool     `yaml:"run"`
	Steps []Action `yaml:"cmds"`
}

type hook struct {
	Name      string   `yaml:"cmdName"`
	Args      []string `yaml:"args"`
	IsSpecial bool     `yaml:"stdoutAsErr,omitempty"`
	Timeout   uint     `yaml:"timeout,omitempty"`
}

// read config file
func LoadCfg(in io.Reader) (Cfg, error) {
	var buf []byte
	c := Cfg{}
	if in == nil || (&in) == nil {
		return c, fmt.Errorf("%w: passed nil as input, please initialize: %q", ErrConfig, in)
	}
	buf, err := io.ReadAll(in)
	if err != nil {
		err = fmt.Errorf("%w: cannot read from %v: %s", ErrConfig, in, err)
		return c, err
	}
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		err = fmt.Errorf("%w: cannot parse yaml file %v: %s", ErrConfig, in, err)
		return c, err
	}
	if len(c) == 0 {
		err = fmt.Errorf("%w: no config data loaded from %v", ErrConfig, in)
		return c, err
	}
	return c, nil
}

// enables writing cfg structure to a file
func newTemplateCfg() Cfg {
	// TODO this function can be refactored (break up)
	const numOfPipes = 2
	build := hook{
		Name: "go",
		Args: []string{"build", "."},
	}
	test := hook{
		Name: "go",
		Args: []string{"test", "."},
	}
	format := hook{
		Name:      "gofmt",
		Args:      []string{"-l", "."},
		IsSpecial: true,
	}
	push := hook{
		Name:    "git",
		Args:    []string{"push", "origin", "main"},
		Timeout: 15,
	}
	m := make(map[string]hook)
	m0 := make(map[string]hook)
	m1 := make(map[string]hook)
	m2 := make(map[string]hook)
	m["build"] = build
	m0["test"] = test
	m2["format"] = format
	//var s []map[string]hook
	var s []Action
	var s1 []Action
	s1 = append(s1, m, m0, m2)
	preCommit := Pipeline{
		Run:   true,
		Steps: s1,
	}
	m1["push"] = push
	s = append(s, m, m0, m2, m1)
	prePush := Pipeline{
		Run:   false,
		Steps: s,
	}
	pipe := make(map[string]Pipeline, numOfPipes)
	pipe["pre-commit"] = preCommit
	pipe["pre-push"] = prePush
	return pipe
}

func MakeTemplateCfg(out io.Writer) error {
	var data []byte
	c := newTemplateCfg()
	data, err := yaml.Marshal(c)
	if err != nil {
		err = fmt.Errorf("%w: cannot marshal to %v: %s", ErrConfig, out, err)
		return err
	}
	N := len(data)
	n, err := out.Write(data)
	if err != nil {
		err = fmt.Errorf("%w: cannot write to %v. Wrote %d, expected: %d bytes: %s", ErrConfig, out, n, N, err)
		return err
	}
	return nil
}

// print to out - mostly for debugging purposes
func printCfg(cfg *Cfg, out io.Writer) {
	for pipename, content := range *cfg {
		fmt.Fprintf(out, "######################\n")
		fmt.Fprintf(out, "found pipename: %q\n- properties:\n - run: %v\n", pipename, content.Run)
		fmt.Fprintf(out, " %q has following steps:\n", pipename)
		for name, vals := range content.Steps {
			fmt.Fprintf(out, "step: %d\n", name)
			for stepName, hook := range vals {
				fmt.Fprintf(out, "  - step name: %s\n  - exe: %q\n  - args: %v\n", stepName, hook.Name, hook.Args)
			}
		}
	}
}
