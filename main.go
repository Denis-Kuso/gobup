/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"
	"fmt"

	"example.com/gobup/cmd"
)

func main() {
	cmd.Execute()
	c := newExampleCfg()
	file, err := os.Create("temp_config.yaml")
	defer file.Close()
	if err != nil {
		fmt.Printf("could not create file: %v\n", err)
		return
	}
	err = makeTemplate(c, file)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return
}
