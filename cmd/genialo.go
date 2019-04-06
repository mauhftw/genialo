package cmd

import (
	// System
	"fmt"
	"os"

	// 3rd Party
	"github.com/spf13/cobra"
)

// set command properties
var rootCmd = &cobra.Command{
	Use: "genialo",
	Long: `Application designed for automating processes and infrastructure
For more info, Please check the repo https://github.com/mauhftw/genialo`,
}

// TODO: Add logs
// TODO: Add comments
func Root() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}
