package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/kainnsoft/migrator/internal/entity"
)

const (
	mySqlTimeout = 30 * time.Second
)

type (
	IMySQLRepo interface {
		GetCountAMoneyMoves(ctx context.Context) (int, error)
		ListAMoneyMoves(ctx context.Context, limit int) (res []*entity.AMoneyMoves, err error)
	}
	mySqlRepo struct {
		db     *sql.DB
		logger *zap.Logger
	}
)

func NewMySQLRepo(db *sql.DB, logger *zap.Logger) IMySQLRepo {
	return &mySqlRepo{
		db:     db,
		logger: logger,
	}
}

var aMoneyMovesFields = []string{
	"id",
	"payment_system",
	"a_id",
	"purse",
	"amount_in_request_currency",
	"from_a_currency",
	"to_a_currency",
	"amount_in_account_currency",
	"exchange_rate",
	"created_dt",
	"accepted_dt",
	"status",
	"comment",
	"fio",
	"ps_request_id",
	"extra_data",
	"reject_reason",
	"reject_reason_text",
	"amount_usd",
	"bank_info",
	"md5",
	"is_user_vip",
	"manager_id",
	"transaction_id",
	"mt_order",
	"custom_cf_account_id",
}

func (m *mySqlRepo) ListAMoneyMoves(ctx context.Context, limit int) (res []*entity.AMoneyMoves, err error) {
	var (
		rows   *sql.Rows
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(ctx, mySqlTimeout)
	defer cancel()

	var query = fmt.Sprintf("SELECT %s FROM cabinet.a_money_moves",
		strings.Join(aMoneyMovesFields, ","))
	if limit > 0 {
		query = query + " LIMIT ?"
		if rows, err = m.db.QueryContext(ctx, query, limit); err != nil {
			return nil, err
		}
	} else {
		if rows, err = m.db.QueryContext(ctx, query); err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var item entity.AMoneyMoves
		if err = rows.Scan(
			&item.ID,
			&item.PaymentSystem,
			&item.AID,
			&item.Purse,
			&item.AmountInRequestCurrency,
			&item.FromACurrency,
			&item.ToACurrency,
			&item.AmountInAccountCurrency,
			&item.ExchangeRate,
			&item.CreatedAt,
			&item.AcceptedAt,
			&item.Status,
			&item.Comment,
			&item.Fio,
			&item.PsRequestID,
			&item.ExtraData,
			&item.RejectReason,
			&item.RejectReasonText,
			&item.AmountUsd,
			&item.BankInfo,
			&item.Md5,
			&item.IsUserVip,
			&item.ManagerID,
			&item.TransactionID,
			&item.MtOrder,
			&item.CustomCfAccountID,
		); err != nil {
			m.logger.Error("err =", zap.Error(err))
		}
		res = append(res, &item)
	}

	return res, nil
}

func (m *mySqlRepo) GetCountAMoneyMoves(ctx context.Context) (int, error) {
	var (
		res    int
		err    error
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(ctx, mySqlTimeout)
	defer cancel()

	if err = m.db.QueryRowContext(ctx, "SELECT count(*) as count FROM cabinet.a_money_moves").Scan(&res); err != nil {
		return 0, err
	}

	return res, nil
}
