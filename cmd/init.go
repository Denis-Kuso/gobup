/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Denis-Kuso/gobup/internal/config"
	"github.com/spf13/cobra"
)

const (
	cfgName = ".gobup.yaml"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init <project_dir>",
	Short:   "Initialize tool in provided directory",
	Long:    `Creates config file named .gobup.yaml in provided directory.`,
	Example: " - create config in current dir\n   gobup init .",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		project := args[0]
		err := createTemplate(project)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
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
