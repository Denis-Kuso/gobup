/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Denis-Kuso/gobup/internal/actions"
	"github.com/Denis-Kuso/gobup/internal/config"
	"github.com/spf13/cobra"
)

const (
	pipelineArg = "pipeline"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run a pipeline (set of commands) specified in .gobup.yaml",
	Example: " - run all pipelines whose 'run' field is set to true\n   gobup run\n - run commands associated with pre-push pipeline\n   gobup run -p pre-push\n",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// ####### args preparation
		var project string
		// assume user means current dir
		if len(args) == 0 {
			project = "."
		} else {
			project = args[0]
		}
		fname, err := filepath.Abs(project)
		if err != nil {
			return fmt.Errorf("%w: project: %s: %v", ErrValidation, project, err)
		}
		pipeline, err := cmd.Flags().GetString(pipelineArg)
		if err != nil {
			return fmt.Errorf("%w: %v: %q", ErrValidation, err, pipeline)
		}
		// ####### args preparation
		cfg, err := os.Open(cfgName)
		defer cfg.Close()
		if err != nil {
			return fmt.Errorf("cannot open config file: %v", err)
		}
		steps, err := preparePipes(cfg, pipeline)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		// any pipelines to run?
		if len(steps) == 0 {
			return fmt.Errorf("nothing to run. Are run fields of your all your pipelines set to false?")
		}
		koraci := makeExeSteps(steps, fname)
		for _, korak := range koraci {
			err := korak.Execute()
			if err != nil {
				fmt.Printf(" step: %-10q --> \033[31mFAILURE\033[0m\n", korak.Name)
				return err
			}
			fmt.Printf(" step: %-10q --> \033[32mSUCCESS\033[0m\n", korak.Name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP(pipelineArg, "p", "", "specific pipeline (e.g. \"pre-commit\") defined in your .gobup.yaml")
}

// if pipeline provided and present in cfg, add to queue
// else add pipelines where run == true to a queue
func preparePipes(cfg io.Reader, pipeline string) ([]config.Action, error) {
	var red []config.Action
	c, err := config.LoadCfg(cfg)
	if err != nil {
		return red, err
	}
	if pipeline != "" {
		pipe, ok := c[pipeline]
		if !ok {
			return red, fmt.Errorf("pipeline: %q not found", pipeline)
		}
		red = append(red, pipe.Steps...)
	} else {
		for _, pipe := range c {
			if pipe.Run {
				red = append(red, pipe.Steps...)
			}
		}
	}
	return red, nil
}

func makeExeSteps(pipelines []config.Action, project string) []actions.Step {
	var steps []actions.Step
	for _, pipe := range pipelines {
		s := pipe
		for name, cmd := range s {
			step := actions.NewStep(name, cmd.Name, cmd.Args, project, time.Duration(cmd.Timeout), cmd.IsSpecial)
			steps = append(steps, step)
		}
	}
	return steps
}
