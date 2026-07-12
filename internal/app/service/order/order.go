package sorder

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	ccatalog "github.com/Lagwick/order-service/internal/app/client/catalog"
	"github.com/Lagwick/order-service/internal/app/entity"
	"github.com/Lagwick/order-service/internal/app/repository"
	"github.com/Lagwick/order-service/internal/app/service"
	catalogv1 "github.com/Lagwick/order-service/internal/pkg/grpc/gen/catalog/v1"
)

type srv struct {
	repoOrder   repository.Order
	catalogGrpc ccatalog.Client
}

func NewService(repoOrder repository.Order, catalogGrpc ccatalog.Client) service.Order {
	return &srv{
		repoOrder:   repoOrder,
		catalogGrpc: catalogGrpc,
	}
}

func (s *srv) Create(ctx context.Context, req entity.RequestOrderCreate) (entity.Order, error) {
	productGUIDs := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		productGUIDs = append(productGUIDs, item.ProductGUID.String())
	}

	productsResp, err := s.catalogGrpc.GetProducts(ctx, &catalogv1.GetProductsRequest{
		Guids: productGUIDs,
	})
	if err != nil {
		return entity.Order{}, err
	}

	if len(productsResp.GetMissingGuids()) > 0 {
		return entity.Order{}, entity.ErrIncorrectParameters
	}

	priceByProductGUID := make(map[uuid.UUID]int64, len(productsResp.GetProducts()))
	for _, product := range productsResp.GetProducts() {
		productGUID, err := uuid.FromString(product.GetGuid())
		if err != nil {
			return entity.Order{}, entity.ErrIncorrectParameters
		}

		priceByProductGUID[productGUID] = product.GetPrice()
	}

	now := time.Now()
	orderGUID := uuid.Must(uuid.NewV4())

	items := make([]entity.OrderItem, 0, len(req.Items))
	var totalPrice int64

	for _, item := range req.Items {
		unitPrice, ok := priceByProductGUID[item.ProductGUID]
		if !ok {
			return entity.Order{}, entity.ErrIncorrectParameters
		}

		totalPrice += int64(item.Quantity) * unitPrice

		items = append(items, entity.OrderItem{
			GUID:        uuid.Must(uuid.NewV4()),
			OrderGUID:   orderGUID,
			ProductGUID: item.ProductGUID,
			Quantity:    item.Quantity,
			UnitPrice:   unitPrice,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	order := entity.Order{
		GUID:       orderGUID,
		UserGUID:   req.UserGUID,
		TotalPrice: totalPrice,
		Currency:   req.Currency,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
		Items:      items,
	}

	if err := s.repoOrder.Create(ctx, order); err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (s *srv) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Order, error) {
	return s.repoOrder.GetByGUID(ctx, guid)
}

func (s *srv) List(ctx context.Context, req entity.RequestOrderList) ([]entity.Order, error) {
	return s.repoOrder.List(ctx, req.Status, req.UserGUID)
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	return s.repoOrder.InsideTx(ctx, func(ctx context.Context) error {
		if _, err := s.repoOrder.GetByGUID(ctx, guid); err != nil {
			return err
		}

		return s.repoOrder.Delete(ctx, guid)
	})
}

func (s *srv) Update(
	ctx context.Context,
	guid uuid.UUID,
	req entity.RequestOrderUpdate,
) (entity.Order, error) {
	var order entity.Order

	err := s.repoOrder.InsideTx(ctx, func(ctx context.Context) error {
		var err error

		order, err = s.repoOrder.GetByGUID(ctx, guid)
		if err != nil {
			return err
		}

		order.Status = req.Status
		order.UpdatedAt = time.Now()

		if err := s.repoOrder.Update(ctx, order); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}
