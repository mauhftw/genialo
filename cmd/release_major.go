package cmd

import (
	//"fmt"

	"github.com/spf13/cobra"
)

// major represents the major command
var majorCmd = &cobra.Command{
	Use:   "major",
	Short: "Creates a new major release",
	Run:   GithubHandler,
}

func init() {
	majorCmd.Flags().StringVarP(&Application, "application", "a", "", "The name of the application you want to release")
	majorCmd.Flags().StringVarP(&Organization, "organization", "o", "SweatWorks", "The name of the organization/owner of the repo")
	majorCmd.MarkFlagRequired("application")
}
