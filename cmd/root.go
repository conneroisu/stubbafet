/*
Copyright © 2024 conneroisu <conneroisu@outlook.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stubba",
	Short: "Generate python stubs for your python code from a requirements file.",
	Long: `Stubba is a tool that generates python stubs for your python code from a requirements file.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd.Context())
		fmt.Println("Stub generation complete.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
