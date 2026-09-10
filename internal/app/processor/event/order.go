package eprocessor

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/Lagwick/order-service/internal/app/entity"
	ehandler "github.com/Lagwick/order-service/internal/app/handler/event"
	"github.com/Lagwick/order-service/internal/app/processor"
	"github.com/Lagwick/order-service/internal/pkg/broker"
)

type orderDeliveryCalculatedProc struct {
	h   ehandler.OrderDeliveryCalculated
	bus broker.Bus[entity.EventOrderDeliveryCalculated]
}

func NewOrderDeliveryCalculatedEventsCatcher(
	h ehandler.OrderDeliveryCalculated,
	bus broker.Bus[entity.EventOrderDeliveryCalculated],
) processor.Processor {
	return &orderDeliveryCalculatedProc{
		h:   h,
		bus: bus,
	}
}

func (p *orderDeliveryCalculatedProc) StartAsync(
	ctx context.Context,
	wg *sync.WaitGroup,
) {
	if err := p.bus.Subscribe(
		ctx,
		wg,
		p.h.CallbackOrderDeliveryCalculated,
	); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to subscribe to order.delivery.calculated events")
	}

	log.Info().
		Msg("order.delivery.calculated subscription requested")
}
