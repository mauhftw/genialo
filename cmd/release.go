package cmd

import (
	"github.com/spf13/cobra"
)

// release flags
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
