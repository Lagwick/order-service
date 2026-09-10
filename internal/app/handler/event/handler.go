package ehandler

import (
	"context"

	"github.com/Lagwick/order-service/internal/app/entity"
	"github.com/Lagwick/order-service/internal/pkg/broker"
)

type OrderDeliveryCalculated interface {
	CallbackOrderDeliveryCalculated(
		ctx context.Context,
		ev *entity.EventOrderDeliveryCalculated,
		headers []broker.Header,
	) error
}
