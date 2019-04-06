package cmd

import (
	// 3rd Party
	"github.com/spf13/cobra"
)

// Release flags
var Application string
var DirectoryBase string
var Organization string

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Creates a new software release following semantic versioning approach",
}

func init() {
	releaseCmd.AddCommand(majorCmd, minorCmd, patchCmd)
}
