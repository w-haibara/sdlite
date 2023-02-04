package cli

import (
	"log"

	"github.com/spf13/cobra"
)

func Exec() {
	if err := cmdRoot().Execute(); err != nil {
		log.Fatal("Failed to execute command:", err)
	}
}

func cmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sdlite <command> [flags]",
		Short: "",
	}

	// cmd.Flags().Bool("version", false, "Show sdlite version")
	// cmd.PersistentFlags().Bool("help", false, "Print help")

	cmd.AddCommand(cmdServe())

	return cmd
}
