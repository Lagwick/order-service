package ccatalog

import (
	"context"

	catalogv1 "github.com/Lagwick/order-service/internal/pkg/grpc/gen/catalog/v1"
)

type Client interface {
	Ping(ctx context.Context) error

	GetProduct(
		ctx context.Context,
		req *catalogv1.GetProductRequest,
	) (*catalogv1.GetProductResponse, error)

	GetProducts(
		ctx context.Context,
		req *catalogv1.GetProductsRequest,
	) (*catalogv1.GetProductsResponse, error)
}
