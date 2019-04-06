package cmd

import (
	// 3rd Party
	"github.com/spf13/cobra"
)

// Represents the minor command
var minorCmd = &cobra.Command{
	Use:   "minor",
	Short: "Creates a new minor release",
	Run:   GithubHandler,
}

func init() {
	minorCmd.Flags().StringVarP(&Application, "application", "a", "", "The name of the application you want to release")
	minorCmd.Flags().StringVarP(&Organization, "organization", "o", "SweatWorks", "The name of the organization/owner of the repo")
	minorCmd.MarkFlagRequired("application")
}
