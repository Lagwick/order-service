package rcpostgres

import (
	"context"

	"gorm.io/gorm"
)

type contextKeyTx struct{}

func getTxFromCtx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextKeyTx{}).(*gorm.DB)
	if !ok {
		return nil
	}

	return tx
}

func ctxWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, contextKeyTx{}, tx)
}
