package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// minor represents the minor command
var minorCmd = &cobra.Command{
	Use:   "minor",
	Short: "Creates a new minor release",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("minor called")
	},
}

func init() {
	minorCmd.Flags().StringVarP(&DirectoryBase, "directory-base", "d", "", "The current directory you will work on")
	minorCmd.Flags().StringVarP(&Application, "application", "a", "", "The name of the application you want to release")
	minorCmd.MarkFlagRequired("application")
	minorCmd.MarkFlagRequired("directory-base")

}
