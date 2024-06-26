/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"

	"github.com/Denis-Kuso/gobup/internal/config"
	"github.com/spf13/cobra"
)

const (
	simulateExec = "dry-run"
	//errAsWarn = "err-as-warn"
	pipeline = "pipeline"
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
		akcija, err := cmd.Flags().GetString(pipeline)
		_, err = cmd.Flags().GetBool(simulateExec)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		steps, err := preparePipes(akcija)
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
		}
		for i, s := range steps {
			fmt.Printf("step %d: :%v\n", i, s)
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

// you already have "run" implemented - now you need to prep the state, such
// given the flags and options provided "run" runs according to state
// How?
// [x] load cfg 
// if pipeline provided and present in cfg, add to queue
  // else add pipelines where run == true to a queue
// if n == true, print queue.commands
// else run queue
func preparePipes(cfg io.Reader, pipeline string) ([]config.Action, error) {
	// perhaps make steps immediateley - for that, need to export a type
	var red []config.Action
	c, err := config.LoadCfg(cfg)
	if err != nil {
		return red, err
	}
	if pipeline != "" {
		fmt.Printf("does %s exist?\n", pipeline)
		pipe, ok := c[pipeline]
		if ! ok{
			return red, fmt.Errorf("%q not found in %v", pipeline, cfg)
		}
		red = append(red, pipe.Steps...)
	}else {
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
func makeExeSteps(pipelines []config.Action) error {
	var steps []string
	for _, pipe := range pipelines {
		s := pipe
		for name, cmd := range s {
			// I should make a "NewStep"/NewTimeoutStep" here
			// Where should NewStep be located?
			steps = append(steps, cmd.Name)
				// how do I pass the project??
				// timeout option?
		}
	}
	return nil
}


// what about warnings?
// I could ignore this feature for know
// I could have a warnings channel, to which I add, instead of the
// errs channel
