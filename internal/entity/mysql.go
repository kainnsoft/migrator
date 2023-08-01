package entity

import (
	"database/sql"
	"time"
)

type AMoneyMoves struct {
	ID                      int64           `db:"id"`
	PaymentSystem           string          `db:"payment_system"`
	AID                     int64           `db:"a_id"`
	Purse                   sql.NullString  `db:"purse"`
	AmountInRequestCurrency sql.NullFloat64 `db:"amount_in_request_currency"`
	FromACurrency           string          `db:"from_a_currency"`
	ToACurrency             string          `db:"to_a_currency"`
	AmountInAccountCurrency sql.NullFloat64 `db:"amount_in_account_currency"`
	ExchangeRate            float64         `db:"exchange_rate"`
	CreatedAt               time.Time       `db:"created_dt"`
	AcceptedAt              sql.NullTime    `db:"accepted_dt"`
	Status                  string          `db:"status"`
	Comment                 sql.NullString  `db:"comment"`
	Fio                     sql.NullString  `db:"fio"`
	PsRequestID             sql.NullString  `db:"ps_request_id"`
	ExtraData               sql.NullString  `db:"extra_data"`
	RejectReason            sql.NullString  `db:"reject_reason"`
	RejectReasonText        sql.NullString  `db:"reject_reason_text"`
	AmountUsd               sql.NullFloat64 `db:"amount_usd"`
	BankInfo                sql.NullString  `db:"bank_info"`
	Md5                     sql.NullString  `db:"md5"`
	IsUserVip               int             `db:"is_user_vip"`
	ManagerID               sql.NullInt64   `db:"manager_id"`
	TransactionID           sql.NullInt64   `db:"transaction_id"`
	MtOrder                 sql.NullInt64   `db:"mt_order"`
	CustomCfAccountID       sql.NullInt64   `db:"custom_cf_account_id"`
}
