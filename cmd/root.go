// Package cmd defines the implementations of the root command
// and the subcommands associated with the gobup tool.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// default values if nothing provided during build/linking phase
	Version    = "development"
	CommitHash = "not provided"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gobup",
	Short: "Local CI/pipeline builder",
	Long: `Configuration based pipeline building. 
	Edit the provided config file and define custom pipelines (e.g. "dev", "pre-commit")
	to match your workflow and/or a policy (e.g. no pushing of untested code). 

	A pipeline can be ran manually whenever you want. You could add the command
	you use to run a pipeline into a desired hook file, such as pre-commit.
	That pipeline will then run every time you try to commit.`,
	Version:      Version,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	versionTemplate := fmt.Sprintf("%s - version %s (commitHash: %s)\n", rootCmd.Use, Version, CommitHash)
	rootCmd.SetVersionTemplate(versionTemplate)
}
