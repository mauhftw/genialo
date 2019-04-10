package cmd

import (
	// 3rd Party
	"github.com/spf13/cobra"
)

// Represents the major command
var majorCmd = &cobra.Command{
	Use:   "major",
	Short: "Creates a new major release",
	Run:   GithubReleaseHandler,
}

func init() {
	majorCmd.Flags().StringVarP(&Application, "application", "a", "", "The name of the application you want to release")
	majorCmd.Flags().StringVarP(&Organization, "organization", "o", "SweatWorks", "The name of the organization/owner of the repo")
	majorCmd.MarkFlagRequired("application")
}
