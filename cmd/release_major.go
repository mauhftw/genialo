package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// major represents the major command
var majorCmd = &cobra.Command{
	Use:   "major",
	Short: "Creates a new major release",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("major called")
	},
}

func init() {
	majorCmd.Flags().StringVarP(&DirectoryBase, "directory-base", "d", "", "The current directory you will work on")
	majorCmd.Flags().StringVarP(&Application, "application", "a", "", "The name of the application you want to release")
	majorCmd.MarkFlagRequired("application")
	majorCmd.MarkFlagRequired("directory-base")
}
