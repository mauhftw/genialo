package cmd

import (
	// 3rd Party
	"github.com/spf13/cobra"
)

// Labels flags
var LabelFile string

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage Github labels",
}

func init() {
	labelsCmd.AddCommand(createLabelsCmd)
	labelsCmd.AddCommand(deleteLabelsCmd)
}
