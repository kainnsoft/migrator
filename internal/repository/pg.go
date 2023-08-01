package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/kainnsoft/migrator/internal/entity"
)

type (
	IPGRepo interface {
		GetCountPayments(ctx context.Context) (int, error)
		ListPayments(ctx context.Context, limit int) (res []*entity.Payment, err error)
		InsertPayment(ctx context.Context, payment *entity.Payment) (rowID int64, err error)
		InsertPaymentWithConflictResolwing(ctx context.Context, payment *entity.Payment) (rowID int64, err error)
	}
	pgRepo struct {
		db     *pgxpool.Pool
		logger *zap.Logger
	}
)

func NewPGRepo(db *pgxpool.Pool, logger *zap.Logger) IPGRepo {
	return &pgRepo{
		db:     db,
		logger: logger,
	}
}

var paymentFields = []string{
	"id",
	"created_at",
	"updated_at",
	"accepted_at",
	"operation_type",
	"code",
	"account_id",
	"amount",
	"amount_currency",
	"account_amount",
	"exchange_rate",
	"status",
	"comment",
	"mt_order_id",
	"transaction_id",
	"extra_data",
	"purse",
}

func (p *pgRepo) InsertPayment(
	ctx context.Context,
	payment *entity.Payment) (rowID int64, err error) {
	var (
		sqlCreatePayment = `INSERT INTO public.payments (%s) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		    RETURNING id;`
		queryStr = fmt.Sprintf(sqlCreatePayment, strings.Join(paymentFields, ","))
	)
	err = p.db.QueryRow(ctx, queryStr,
		payment.ID,
		payment.CreatedAt,
		payment.UpdatedAt,
		payment.AcceptedAt,
		payment.OperationType,
		payment.Code,
		payment.AccountID,
		payment.Amount,
		payment.AmountCurrency,
		payment.AccountAmount,
		payment.ExchangeRate,
		payment.Status,
		payment.Comment,
		payment.MtOrderID.Int64,
		payment.TransactionID,
		payment.ExtraData,
		payment.Purse,
	).Scan(&rowID)

	return rowID, err
}

func (p *pgRepo) ListPayments(ctx context.Context, limit int) (res []*entity.Payment, err error) {
	var (
		rows   pgx.Rows
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(ctx, mySqlTimeout)
	defer cancel()

	var query = fmt.Sprintf("SELECT %s FROM public.payments",
		strings.Join(paymentFields, ","))
	if limit > 0 {
		query = query + " LIMIT ?"
		if rows, err = p.db.Query(ctx, query, limit); err != nil {
			return nil, err
		}
	} else {
		if rows, err = p.db.Query(ctx, query); err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var item entity.Payment
		if err = rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.AcceptedAt,
			&item.OperationType,
			&item.Code,
			&item.AccountID,
			&item.Amount,
			&item.AmountCurrency,
			&item.AccountAmount,
			&item.ExchangeRate,
			&item.Status,
			&item.Comment,
			&item.MtOrderID,
			&item.TransactionID,
			&item.ExtraData,
			&item.Purse,
		); err != nil {
			p.logger.Error("err =", zap.Error(err))
		}
		res = append(res, &item)
	}

	return res, nil
}

func (p *pgRepo) GetCountPayments(ctx context.Context) (int, error) {
	var (
		res    int
		err    error
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(ctx, mySqlTimeout)
	defer cancel()

	if err = p.db.QueryRow(ctx, "SELECT count(*) as count FROM public.payments").Scan(&res); err != nil {
		return 0, err
	}

	return res, nil
}

func (p *pgRepo) InsertPaymentWithConflictResolwing(
	ctx context.Context,
	payment *entity.Payment) (rowID int64, err error) {
	var (
		sqlCreatePayment = `INSERT INTO public.payments (%s) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		ON CONFLICT (id) DO UPDATE 
		    SET id = EXCLUDED.id,
		     	created_at = EXCLUDED.created_at,
				updated_at = EXCLUDED.updated_at,
				accepted_at = EXCLUDED.accepted_at,
				operation_type = EXCLUDED.operation_type,
				code = EXCLUDED.code,
				account_id = EXCLUDED.account_id,
				amount = EXCLUDED.amount,
				amount_currency = EXCLUDED.amount_currency,
				account_amount = EXCLUDED.account_amount,
				exchange_rate = EXCLUDED.exchange_rate,
				status = EXCLUDED.status,
				comment = EXCLUDED.comment,
				mt_order_id = EXCLUDED.mt_order_id,
				transaction_id = EXCLUDED.transaction_id,
				extra_data = EXCLUDED.extra_data,
				purse = EXCLUDED.purse 
		    RETURNING id;`
		queryStr = fmt.Sprintf(sqlCreatePayment, strings.Join(paymentFields, ","))
	)
	err = p.db.QueryRow(ctx, queryStr,
		payment.ID,
		payment.CreatedAt,
		payment.UpdatedAt,
		payment.AcceptedAt,
		payment.OperationType,
		payment.Code,
		payment.AccountID,
		payment.Amount,
		payment.AmountCurrency,
		payment.AccountAmount,
		payment.ExchangeRate,
		payment.Status,
		payment.Comment,
		payment.MtOrderID.Int64,
		payment.TransactionID,
		payment.ExtraData,
		payment.Purse,
	).Scan(&rowID)

	return rowID, err
}
