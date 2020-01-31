package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewDefaultCommand() *cobra.Command {
	var defaultCmd = &cobra.Command{
		Use:   "default",
		Short: "Prints out the default internal config and template",
		Long:  "Prints out the default internal config and template you can use as a starting point for developing your own fake data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Printing the default config:")
			fmt.Println("===========================\n")
			fmt.Println(gDefaultConfig)
			fmt.Println("\n\nPrinting the default template:")
			fmt.Println("===========================\n")
			fmt.Println(gDefaultTemplate)
			return nil
		},
	}

	return defaultCmd
}
