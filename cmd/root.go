/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "paralix",
	Short: "A brief description of your application",
	Long: `Paralix is a command-line interface (CLI) tool that simplifies the process
of parallelizing tasks. When working with large sets of data or complex calculations,
executing commands in parallel can significantly reduce the amount
of time required to complete the tasks. `,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
