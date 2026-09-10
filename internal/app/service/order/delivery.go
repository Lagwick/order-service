package sorder

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/Lagwick/order-service/internal/app/repository"
	"github.com/Lagwick/order-service/internal/app/service"
)

type deliverySrv struct {
	repoOrder repository.Order
}

func NewDeliveryService(repoOrder repository.Order) service.OrderDelivery {
	return &deliverySrv{
		repoOrder: repoOrder,
	}
}

func (s *deliverySrv) ApplyDeliveryCalculated(
	ctx context.Context,
	guid uuid.UUID,
	deliveryPrice int64,
) (bool, error) {
	return s.repoOrder.ApplyDeliveryCalculated(
		ctx,
		guid,
		deliveryPrice,
	)
}
