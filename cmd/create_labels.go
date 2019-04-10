package cmd

import (
	// 3rd Party
	"github.com/spf13/cobra"
)

// Represents the create labels command
var createLabelsCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new set of github labels",
	Run:   GithubLabelCreatorHandler,
}

func init() {
	createLabelsCmd.Flags().StringVarP(&Application, "application", "a", "", "The name of the application you want to release")
	createLabelsCmd.Flags().StringVarP(&Organization, "organization", "o", "SweatWorks", "The name of the organization/owner of the repo")
	createLabelsCmd.Flags().StringVarP(&LabelFile, "file", "f", "", "The json file of the labels")
	createLabelsCmd.MarkFlagRequired("application")
	createLabelsCmd.MarkFlagRequired("file")
}
