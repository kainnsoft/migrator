package handler

import (
	"context"

	"go.uber.org/zap"

	v1 "github.com/kainnsoft/migrator/internal/controller/http/v1"
	"github.com/kainnsoft/migrator/internal/entity"
	"github.com/kainnsoft/migrator/internal/usecase"
)

type httpHandler struct {
	mySQLUsecase usecase.IMySQLUsecase
	pgUsecase    usecase.IPGUsecase
	logger       *zap.Logger
}

func NewHandler(
	mySQLUsecase usecase.IMySQLUsecase,
	pgUsecase usecase.IPGUsecase,
	logger *zap.Logger,
) v1.IHttpHandlers {
	return &httpHandler{
		mySQLUsecase: mySQLUsecase,
		pgUsecase:    pgUsecase,
		logger:       logger,
	}
}

func (h *httpHandler) GetCountAMoneyMoves(ctx context.Context) (int, error) {
	return h.mySQLUsecase.UseCaseGetCountAMoneyMoves(ctx)
}

func (h *httpHandler) ListAMoneyMoves(ctx context.Context, limit int) (res []*entity.AMoneyMoves, err error) {
	return h.mySQLUsecase.UseCaseListAMoneyMoves(ctx, limit)
}

func (h *httpHandler) GetCountPayments(ctx context.Context) (int, error) {
	return h.pgUsecase.UseCaseGetCountPayments(ctx)
}

func (h *httpHandler) DoTransferFromAMoneyMovesToAccountPayments(ctx context.Context, limit int) error {
	var (
		aMoneyMovesList []*entity.AMoneyMoves
		err             error
	)
	if aMoneyMovesList, err = h.mySQLUsecase.UseCaseListAMoneyMoves(ctx, limit); err != nil {
		h.logger.Error("transfer can't get list from mysql:", zap.Error(err))
		return err
	}

	if err = h.mySQLUsecase.TransferFromAMoneyMoves(ctx, aMoneyMovesList); err != nil {
		h.logger.Error("transfer can't do transfer from mysql:", zap.Error(err))
		return err
	}

	return nil
}
