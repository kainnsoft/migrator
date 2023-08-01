package entity

import (
	"database/sql"
	"time"
)

type Payment struct {
	ID             int64         `db:"id"`
	CreatedAt      time.Time     `db:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at"`
	AcceptedAt     sql.NullTime  `db:"accepted_at"`
	OperationType  string        `db:"operation_type"`
	Code           string        `db:"code"`
	AccountID      int64         `db:"account_id"`
	Amount         float64       `db:"amount"`
	AmountCurrency string        `db:"amount_currency"`
	AccountAmount  float64       `db:"account_amount"`
	ExchangeRate   float64       `db:"exchange_rate"`
	Status         string        `db:"status"`
	Comment        string        `db:"comment"`
	MtOrderID      sql.NullInt64 `db:"mt_order_id"`
	TransactionID  sql.NullInt64 `db:"transaction_id"`
	ExtraData      string        `db:"extra_data"`
	Purse          string        `db:"purse"`
}

type Account struct {
	ID                int64  `db:"id"`
	UserID            int64  `db:"user_id"`
	Login             int64  `db:"login"`
	ServerID          int64  `db:"server_id"`
	PartnerLogin      int64  `db:"partner_login"`
	BalanceMultiplier int64  `db:"balance_multiplier"`
	Currency          string `db:"currency"`
	Tariff            string `db:"tariff"`
	IsDemo            bool   `db:"is_demo"`
	Status            string `db:"status"`
	Balance           int64  `db:"balance"`
	Platform          string `db:"platform"`
}

type PaymentSystem struct {
	Code             string
	ParentCode       string
	Title            string
	Group            string
	Currencies       []string
	ExecutionTime    string
	CommissionString string
	WarningMessage   string
	URL              string
	Logo             string
	Description      string
	IsCardsSupported bool
}
