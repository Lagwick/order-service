package cmd

import (
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/Lagwick/order-service/internal/app/builder"
)

func SubscribeDeliveryCalculated() *cli.Command {
	return &cli.Command{
		Name:    "subscribe-delivery-calculated",
		Aliases: []string{"sdc"},
		Usage:   "Subscribe to order.delivery.calculated events",
		Description: strings.TrimSpace(`
Initializes dependencies required for consuming order.delivery.calculated events
and updating orders in the database. Graceful shutdown on SIGINT/SIGTERM.
`),
		Action:          cmdSubscribeDeliveryCalculated,
		HideHelpCommand: true,
	}
}

func cmdSubscribeDeliveryCalculated(cCtx *cli.Context) error {
	b := builder.NewBuilder(cCtx)

	b.BuildConfig()
	b.BuildMonitorOpenTelemetry()

	b.BuildRepoConnPostgres()
	b.BuildRepoOrder()

	b.BuildBrokerKafkaClient()
	b.BuildBusOrderDeliveryCalculated()

	b.BuildServiceOrderDelivery()
	b.BuildHandlerOrderDeliveryCalculated()
	b.BuildProcOrderDeliveryCalculated()

	b.BuildProcHttp()

	b.Run()

	return nil
}
