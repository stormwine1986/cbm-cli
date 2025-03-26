/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var projectName string

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Useful Project Admin Tools",
	Long:  "Project command provides tools for managing projects",
	Run: func(cmd *cobra.Command, args []string) {
		if projectName != "" {
			cmd.Printf("Working with project: %s\n", projectName)
		}
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
