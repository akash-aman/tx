/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tx",
	Short: "CLI tool for generating templates",
	Long:  `tx is a powerful CLI tool designed for creating starter project templates.`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(installCmd, addCmd, rmCmd)
}
