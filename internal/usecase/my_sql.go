package usecase

import (
	"context"
	"encoding/json"

	"github.com/kainnsoft/migrator/internal/repository"

	"github.com/kainnsoft/migrator/internal/entity"
)

type (
	IMySQLUsecase interface {
		UseCaseGetCountAMoneyMoves(context.Context) (int, error)
		UseCaseListAMoneyMoves(ctx context.Context, limit int) (res []*entity.AMoneyMoves, err error)
		TransferFromAMoneyMoves(ctx context.Context, aMoneyMovesList []*entity.AMoneyMoves) (err error)
		TransferFromAMoneyMovesArray(ctx context.Context, aMoneyMovesList []*entity.AMoneyMoves) ([]*entity.Payment, error)
	}
	mySQLUsecase struct {
		Repo    repository.IMySQLRepo
		rmqRepo repository.IRmqRepo
	}
)

func NewMySQLUsecase(repo repository.IMySQLRepo, rmqRepo repository.IRmqRepo) IMySQLUsecase {
	return &mySQLUsecase{
		Repo:    repo,
		rmqRepo: rmqRepo,
	}
}

func (u *mySQLUsecase) UseCaseGetCountAMoneyMoves(ctx context.Context) (int, error) {
	return u.Repo.GetCountAMoneyMoves(ctx)
}

func (u *mySQLUsecase) UseCaseListAMoneyMoves(ctx context.Context, limit int) ([]*entity.AMoneyMoves, error) {
	return u.Repo.ListAMoneyMoves(ctx, limit)
}

func (u *mySQLUsecase) TransferFromAMoneyMoves(
	ctx context.Context,
	aMoneyMovesList []*entity.AMoneyMoves,
) (err error) {
	var (
		payment           = entity.Payment{}
		vExtraData string = "{}"
		data       []byte
	)
	for _, v := range aMoneyMovesList {
		payment.ID = v.ID
		payment.CreatedAt = v.CreatedAt
		// payment.UpdatedAt = v.
		payment.AcceptedAt = v.AcceptedAt
		payment.OperationType = "award" // maybe first word of comment
		payment.Code = v.PaymentSystem
		payment.AccountID = v.AID
		payment.Amount = v.AmountInRequestCurrency.Float64
		payment.AmountCurrency = v.ToACurrency
		payment.AccountAmount = v.AmountInAccountCurrency.Float64
		payment.ExchangeRate = v.ExchangeRate
		payment.Status = v.Status
		payment.Comment = v.Comment.String
		payment.MtOrderID = v.MtOrder
		payment.TransactionID = v.TransactionID
		if v.ExtraData.String != "" {
			payment.ExtraData = v.ExtraData.String
		} else {
			payment.ExtraData = vExtraData
		}
		payment.Purse = v.Purse.String

		if data, err = json.Marshal(payment); err != nil {
			return err
		}
		if err = u.rmqRepo.PushMsg(ctx, data); err != nil {
			return err
		}
	}

	return nil
}

func (u *mySQLUsecase) PushMsg() (err error) {

	return nil
}

func (u *mySQLUsecase) TransferFromAMoneyMovesArray(
	ctx context.Context,
	aMoneyMovesList []*entity.AMoneyMoves,
) ([]*entity.Payment, error) {
	var (
		payment             = entity.Payment{}
		paymentsList        = make([]*entity.Payment, 0, len(aMoneyMovesList))
		vExtraData   string = "{}"
	)
	for _, v := range aMoneyMovesList {
		payment.ID = v.ID
		payment.CreatedAt = v.CreatedAt
		// payment.UpdatedAt = v.
		payment.AcceptedAt = v.AcceptedAt
		payment.OperationType = "award" // maybe first word of comment
		payment.Code = v.PaymentSystem
		payment.AccountID = v.AID
		payment.Amount = v.AmountInRequestCurrency.Float64
		payment.AmountCurrency = v.ToACurrency
		payment.AccountAmount = v.AmountInAccountCurrency.Float64
		payment.ExchangeRate = v.ExchangeRate
		payment.Status = v.Status
		payment.Comment = v.Comment.String
		payment.MtOrderID = v.MtOrder
		payment.TransactionID = v.TransactionID
		if v.ExtraData.String != "" {
			payment.ExtraData = v.ExtraData.String
		} else {
			payment.ExtraData = vExtraData
		}
		payment.Purse = v.Purse.String

		paymentsList = append(paymentsList, &payment)
	}

	return paymentsList, nil
}
