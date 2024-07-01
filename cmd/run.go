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
	simulateExec = "dry-run"
	pipeline     = "pipeline"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a set of commands defined in a pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

to quickly create a Cobra application.`,
	Example: ` gobup run --dry-run -p pre-push -> print commands associated
	with pre-push pipeline
	 gobup run . -> run all pipelines (if their "run" field is set to true.`,
	Run: func(cmd *cobra.Command, args []string) {
		var project string
		if len(args) == 0 {
			project = "."
		} else {
			project = args[0]
		}
		fname, err := filepath.Abs(project)
		if err != nil {
			fmt.Printf("invalid pathname: %s", project)
			return
		}
		akcija, err := cmd.Flags().GetString(pipeline)
		_, err = cmd.Flags().GetBool(simulateExec)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		cfg, err := os.Open(cfgName)
		defer cfg.Close()
		if err != nil {
			fmt.Printf("cannot open config file: %v\n", err)
			return
		}
		steps, err := preparePipes(cfg, akcija)
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
			return
		}
		// any pipelines to run?
		if len(steps) == 0 {
			fmt.Printf("nothing to run\n")
			return
		}
		koraci := makeExeSteps(steps, fname)
		for _, korak := range koraci {
			msg, err := korak.Execute()
			if err != nil {
				fmt.Print(err)
				return
			}
			fmt.Print(msg)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().BoolP(simulateExec, "n", false, "print the commands that would be executed")
	// WILL NOT IMPLEMENT FOR NOW
	//runCmd.Flags().BoolP(errAsWarn, "w", false, "if possible, treat errors as warnings (do not stop execution in a pipeline)")
	//runCmd.MarkFlagsMutuallyExclusive(simulateExec, errAsWarn)
	runCmd.Flags().StringP(pipeline, "p", "", "a specific pipeline (e.g. \"pre-commit\") defined in your config.yaml to run")
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

// perhaps, once/if the pipeline options are irelevant, return/operate on
// "steps" instead...
func makeExeSteps(pipelines []config.Action, project string) []actions.Executer {
	var steps []actions.Executer
	for _, pipe := range pipelines {
		s := pipe
		for name, cmd := range s {
			msg := fmt.Sprintf(" step: %s -> SUCCESS", name)
			step := actions.NewTimeoutStep(name, cmd.Name, cmd.Args, msg, project, time.Duration(cmd.Timeout), cmd.IsSpecial)
			steps = append(steps, step)
		}
	}
	return steps
}
