package providers

import (
	"context"

	"github.com/urfave/cli/v3"
)

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "provider",
			Usage: "a provider to search",
		},
		&cli.StringFlag{
			Name:  "version",
			Usage: "what version of the provide to search",
		},
		&cli.BoolFlag{
			Name:  "data",
			Usage: "whether to search for a data resource instead",
		},
		&cli.StringFlag{
			Name:  "resource",
			Usage: "terraform resource name to search for",
		},
	}
}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "provider",
		Usage: "get provider documentation",
		Flags: flags(),
		Commands: []*cli.Command{
			{
				Name:  "get-doc",
				Usage: "gets documentation for a specific resource",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					provider := cmd.String("provider")
					version := cmd.String("version")
					isData := cmd.Bool("data")
					resource := cmd.String("resource")

					command(provider, version, resource, isData)
					return nil
				},
			},
		},
	}
}
