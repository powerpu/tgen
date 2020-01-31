package main

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "tgen",
	Short: "tgen is a fake time series data generator",
	Long: `
(T)ime series data (Gen)erator is a flexible fake time series data
generator built with love by me, myself and I in Go. Complete documentation is
available at https://github.com/powerpu/tgen
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	rootCmd.AddCommand(NewGenerateCmd())
	rootCmd.AddCommand(NewPlayareaCmd())
	rootCmd.AddCommand(NewDefaultCommand())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
