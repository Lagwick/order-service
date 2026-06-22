package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Lagwick/order-service/cmd"
)

func main() {
	app := &cli.App{
		Name:    "order-service",
		Version: "1.0.0",
		Usage:   "Order management service",
		Commands: []*cli.Command{
			cmd.WebServer(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-json",
				Usage: "Enable console logger instead of JSON",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
