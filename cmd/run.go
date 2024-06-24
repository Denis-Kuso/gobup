/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	simulateExec = "dry-run"
	errAsWarn = "err-as-warn"
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
		_, err := cmd.Flags().GetString(pipeline)
		_, err = cmd.Flags().GetBool(simulateExec)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		_, err = cmd.Flags().GetBool(errAsWarn)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().BoolP(simulateExec, "n", false, "print the commands that would be executed")
	runCmd.Flags().BoolP(errAsWarn, "w", false, "if possible, treat errors as warnings (do not stop execution in a pipeline)")
	runCmd.MarkFlagsMutuallyExclusive(simulateExec, errAsWarn)
	runCmd.Flags().StringP(pipeline, "p", "", "a specific pipeline (e.g. \"pre-commit\") defined in your config.yaml to run")
}
