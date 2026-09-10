package eorder

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"

	"github.com/Lagwick/order-service/internal/app/entity"
	ehandler "github.com/Lagwick/order-service/internal/app/handler/event"
	"github.com/Lagwick/order-service/internal/app/service"
	"github.com/Lagwick/order-service/internal/pkg/broker"
)

type handler struct {
	delivery service.OrderDelivery
}

func NewHandler(
	delivery service.OrderDelivery,
) ehandler.OrderDeliveryCalculated {
	return &handler{
		delivery: delivery,
	}
}

func (h *handler) CallbackOrderDeliveryCalculated(
	ctx context.Context,
	ev *entity.EventOrderDeliveryCalculated,
	_ []broker.Header,
) error {
	log.Info().
		Ctx(ctx).
		Str("order_guid", ev.OrderGUID).
		Int64("delivery_price", ev.DeliveryPrice).
		Str("currency", ev.Currency).
		Msg("Получено событие order.delivery.calculated")

	if ev.OrderGUID == "" {
		return broker.NotCriticalError(
			fmt.Errorf("order.delivery.calculated: empty order_guid"),
		)
	}

	orderGUID, err := uuid.FromString(ev.OrderGUID)
	if err != nil {
		return broker.NotCriticalError(
			fmt.Errorf("order.delivery.calculated: invalid order_guid: %w", err),
		)
	}

	applied, err := h.delivery.ApplyDeliveryCalculated(
		ctx,
		orderGUID,
		ev.DeliveryPrice,
	)
	if err != nil {
		return fmt.Errorf("apply delivery calculated: %w", err)
	}

	if !applied {
		log.Warn().
			Ctx(ctx).
			Str("order_guid", ev.OrderGUID).
			Msg("delivery calculated event skipped")

		return nil
	}

	log.Info().
		Ctx(ctx).
		Str("order_guid", ev.OrderGUID).
		Int64("delivery_price", ev.DeliveryPrice).
		Str("currency", ev.Currency).
		Msg("Заказ обновлён после расчёта доставки")

	return nil
}
