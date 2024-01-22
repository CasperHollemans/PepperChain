package cmd

import (
	"fmt"
	"github.com/CasperHollemans/PepperChain/internal/api"
	"github.com/spf13/cobra"
)

func NewServerCommand() *cobra.Command {
	var (
		importCmd = &cobra.Command{
			Use:   "server",
			Short: "Starts a server",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Starting server...")
				api.StartServer()
			},
		}
	)

	return importCmd
}
