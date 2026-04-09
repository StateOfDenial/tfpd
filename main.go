package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/StateOfDenial/tfpd/cmd/providers"
)

func main() {
	app := &cli.Command{
		Name:  "tf-state-server",
		Usage: "Terraform HTTP state server",
		Commands: []*cli.Command{
			providers.Command(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		// urfave/cli handles error output
		os.Exit(1)
	}
}
