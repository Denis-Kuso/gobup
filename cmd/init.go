/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	noGit = "no-git"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <project_dir>",
	Short: "initialize tool in provided dir",
	Long: `Creates config file in specified <project_dir>. If non-existing, creates
	pre-commit and pre-push files in .git/hooks. For example:

They are then modified such that your whatever you specifed in config.yaml
	will be ran on pre-commit and pre-push.`,
	Example: " gobup init .\n gobup init --no-git ./someProjectDir",
	Args: cobra.ExactArgs(1) ,
	Run: func(cmd *cobra.Command, args []string) {
		noGit, err := cmd.Flags().GetBool(noGit)
		if err != nil {
			fmt.Printf("err: %v", err)
		}
		if !noGit {
			fmt.Printf("normal init, flag: %v\n", noGit)
			fmt.Printf("creating files...\n")
		}else
		{
			fmt.Printf("Don't want to use Git, flag: %v\n", noGit)
		}
		fmt.Printf("Creating a template config.yaml\n")
		fmt.Printf("init called on project %q\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP(noGit, "g", false, "initialize without relying on git")

	// Here you will define your flags and configuration settings.
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
