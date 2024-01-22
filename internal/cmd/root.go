package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "pepperchain",
		Short: "PepperChain, a simplistic blockchain implementation in Go",
	}
	rootCmd.AddCommand(NewServerCommand())
	return rootCmd
}

func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
