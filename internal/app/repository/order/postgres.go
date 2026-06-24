package porder

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"github.com/Lagwick/order-service/internal/app/entity"
	"github.com/Lagwick/order-service/internal/app/repository"
	rcpostgres "github.com/Lagwick/order-service/internal/app/repository/conn/postgres"
)

type repoPg struct {
	conn *rcpostgres.Client
}

func NewRepo(client *rcpostgres.Client) repository.Order {
	return &repoPg{conn: client}
}

func (r *repoPg) InsideTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.conn.InsideTx(ctx, fn)
}

func (r *repoPg) Create(ctx context.Context, order entity.Order) error {
	db := r.conn.GetDB(ctx)

	if err := db.Create(&order).Error; err != nil {
		return err
	}
	return nil
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Order, error) {
	db := r.conn.GetDB(ctx)

	var order entity.Order
	err := db.
		Preload("Items").
		Where("guid = ?", guid).
		First(&order).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Order{}, entity.ErrNotFound
		}
		return entity.Order{}, err
	}
	return order, nil
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	db := r.conn.GetDB(ctx)

	result := db.Where("guid = ?", guid).Delete(&entity.Order{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (r *repoPg) Update(ctx context.Context, order entity.Order) error {
	db := r.conn.GetDB(ctx)

	result := db.Model(&entity.Order{}).Where("guid = ?", order.GUID).Updates(map[string]any{"status": order.Status})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (r *repoPg) List(ctx context.Context, status *string, userGUID *uuid.UUID) ([]entity.Order, error) {
	db := r.conn.GetDB(ctx)

	query := db.Model(&entity.Order{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if userGUID != nil {
		query = query.Where("user_guid = ?", *userGUID)
	}

	var orders []entity.Order
	result := query.Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

// TODO: Реализуйте Update, Delete, List.
