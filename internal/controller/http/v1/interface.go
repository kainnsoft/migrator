package v1

import (
	"context"

	"github.com/kainnsoft/migrator/internal/entity"
)

type IHttpHandlers interface {
	GetCountAMoneyMoves(context.Context) (int, error)
	ListAMoneyMoves(ctx context.Context, limit int) (res []*entity.AMoneyMoves, err error)
	GetCountPayments(context.Context) (int, error)
	DoTransferFromAMoneyMovesToAccountPayments(ctx context.Context, limit int) error
}
