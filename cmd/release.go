package cmd

import (
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Creates a new software release following semantic versioning approach",
}

func init() {
	releaseCmd.AddCommand(majorCmd, minorCmd, patchCmd)
}
