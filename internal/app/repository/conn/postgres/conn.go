package rcpostgres

import (
	"context"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Lagwick/order-service/internal/app/config/section"
)

type Client struct {
	db  *gorm.DB
	cfg section.RepositoryPostgres
}

func (c *Client) DB() *gorm.DB {
	return c.db
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	u := url.URL{
		Scheme: "postgres",
		Host:   cfg.Address,
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Path:   cfg.Name,
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	dsn := u.String()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, err
	}
	return &Client{db: db, cfg: cfg}, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *Client) GetDB(ctx context.Context) *gorm.DB {
	tx := getTxFromCtx(ctx)
	if tx != nil {
		return tx
	}
	return c.db
}

func (c *Client) InsideTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := getTxFromCtx(ctx)
	if tx != nil {
		return fn(ctx)
	}

	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctxWithTx(ctx, tx))
	})
}
