package providers

import (
	"context"

	"github.com/urfave/cli/v3"
)

func flags() []cli.Flag {
	return []cli.Flag{}
}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "provider",
		Usage: "get provider documentation",
		Flags: flags(),
		Commands: []*cli.Command{
			{
				Name:  "get-resource-doc",
				Usage: "gets documentation for a specific resource",
				Action: func(ctx context.Context, c *cli.Command) error {
					command()
					return nil
				},
			},
		},
	}
}
