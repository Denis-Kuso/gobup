/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Denis-Kuso/gobup/internal/config"
	"github.com/spf13/cobra"
)

const (
	noGit             = "no-git"
	cfgName           = ".gobup.yaml"
	preCommitFilename = "pre-commit"
	prePushFilename   = "pre-push"
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
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		project := args[0]
		noGit, err := cmd.Flags().GetBool(noGit)
		if err != nil {
			fmt.Printf("err: %v", err)
		}
		if !noGit {
			fmt.Printf("normal init, flag: %v\n", noGit)
			makeHookFiles(project)
		} else {
			fmt.Printf("Don't want to use Git, flag: %v\n", noGit)
		}
		fmt.Printf("init called on project %q\n", project)
		n, err := filepath.Abs(filepath.Join(project, ".git", preCommitFilename))
		fmt.Printf("filepath: %q, err : %\n", n, err)
		createTemplate(project)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP(noGit, "g", false, "initialize without relying on git")

}

// checks of file existence could be done by another function (overlap of
// functionality with the operations in createTemplate

// generate git hook files from project provided
// generates commands to write into hook files
func makeHookFiles(project string) error {
	//	fname, err := filepath.Abs(filepath.Join(project, ".git", preCommitFilename))
	//	name, err := filepath.Abs(filepath.Join(project, ".git", prePushFilename))
	// commands to write in git hook files
	// not yet complete
	// commands to be ran - currently placeholder until run is implemented
	commit := "gobup run . -a pre-commit"
	push := "gobup run . -a pre-push -e"
	err1 := writeCmd([]byte(commit), os.Stdout)
	err2 := writeCmd([]byte(push), os.Stdout)
	if err := errors.Join(err1, err2); err != nil {
		return err
	}
	return nil
}

func writeCmd(cmd []byte, out io.Writer) error {
	N := len(cmd)
	if N == 0 {
		err := fmt.Errorf("No data: %q to write", string(cmd))
		return err
	}
	n, err := out.Write([]byte(cmd))
	if err != nil {
		msg := fmt.Sprintf("unsuccessful write: %v, wrote %d bytes, expected: %d", err, n, N)
		return fmt.Errorf("%s: %v", msg, err)
	}
	return nil
}

// createTemplate writes template cfg to location of project
// if project not a valid pathname or file with default cfg name
// exists it returns an error
func createTemplate(project string) error {
	// make/check valid path
	fname, err := filepath.Abs(filepath.Join(project, cfgName))
	if err != nil {
		return fmt.Errorf("invalid pathname: %s", project)
	}
	// don't overwrite existing config file
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	defer f.Close()
	if err != nil {
		// don't offer options as of yet (whether to replace)
		err = fmt.Errorf("cannot create template: %v", err)
		return err
	}
	err = config.MakeTemplateCfg(f)
	return err
}
