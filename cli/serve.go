package cli

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/w-haibara/sdlite/server"
)

func cmdServe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve sdlite server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe()
		},
	}

	return cmd
}

func runServe() error {
	if err := server.ListenAndServe(":8080"); err != nil {
		log.Fatal(err)
	}

	return nil
}
