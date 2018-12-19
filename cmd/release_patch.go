package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// patch represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Creates a new patch release",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("patch called")
	},
}

func init() {
	patchCmd.Flags().StringVarP(&DirectoryBase, "directory-base", "d", "", "The current directory you will work on")
	patchCmd.Flags().StringVarP(&Application, "application", "a", "", "The name of the application you want to release")
	patchCmd.MarkFlagRequired("application")
	patchCmd.MarkFlagRequired("directory-base")
}
