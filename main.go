package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/StateOfDenial/tfpd/cmd/providers"
)

func main() {
	app := &cli.Command{
		Name:  "tfpd",
		Usage: "Terraform provider docs getter",
		Commands: []*cli.Command{
			providers.Command(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		// urfave/cli handles error output
		os.Exit(1)
	}
}
