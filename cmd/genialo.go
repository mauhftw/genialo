package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// FLAGS
// this should be defined in other place
var Application string
var DirectoryBase string

// set command properties
var rootCmd = &cobra.Command{
	Use: "genialo",
	Long: `Application designed for automating processes and infrastructure
For more info, Please check the repo https://github.com/mauhftw/genialo`,
}

// agregar logs
// agregar comentarios
func Root() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}
