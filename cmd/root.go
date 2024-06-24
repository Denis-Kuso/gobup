/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gobup",
	Short: "a local CI/pipeline builder",
	Long: `Name a set of steps in a yaml file. Then those steps can be run 
for you like a CI pipeline. This be can done as part of commits or pushes. Or
	manually whenever you please.`,
	Version: "dev",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gobup.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	versionTemplate := fmt.Sprintf("%s: %s - version %q (commitHash?)\n", rootCmd.Name(), rootCmd.Short, rootCmd.Version)
	rootCmd.SetVersionTemplate(versionTemplate)

}


