package cmd

import (
	// 3rd Party
	"github.com/spf13/cobra"
)

// Represents the delete labels command
var deleteLabelsCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes all github labels",
	Run:   GithubLabelDestroyerHandler,
}

func init() {
	deleteLabelsCmd.Flags().StringVarP(&Application, "application", "a", "", "The name of the application you want to release")
	deleteLabelsCmd.Flags().StringVarP(&Organization, "organization", "o", "SweatWorks", "The name of the organization/owner of the repo")
	deleteLabelsCmd.MarkFlagRequired("application")
}
