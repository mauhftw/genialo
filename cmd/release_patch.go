package cmd

import (
	// 3rd Party
	"github.com/spf13/cobra"
)

// Represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Creates a new patch release",
	Run:   GithubHandler,
}

func init() {
	patchCmd.Flags().StringVarP(&Application, "application", "a", "", "The name of the application you want to release")
	patchCmd.Flags().StringVarP(&Organization, "organization", "o", "SweatWorks", "The name of the organization/owner of the repo")
	patchCmd.MarkFlagRequired("application")
}
