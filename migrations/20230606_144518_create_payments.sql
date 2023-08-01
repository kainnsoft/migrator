-- +goose Up

CREATE TYPE PAYMENT_TYPE AS ENUM ('deposit', 'withdrawal', 'award', 'in', 'pin', 'ex', 'rollback', 'bonus', 'payment', 'dividend');
CREATE TYPE OPERATION_STATUS AS ENUM ('new', 'completed', 'error', 'accepted', 'rejected', 'inProgress', 'hasquestion', 'rollback');

-- Основная таблица, куда пишем вместо MySQL
CREATE TABLE IF NOT EXISTS public.payments
(
    id               BIGSERIAL primary key,                            -- ID пеймента (a_money_moves.id (from MySQL))
    created_at       TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT now(), -- Дата создания платежа (a_money_moves.created_dt → created_at (from MySQL))
    updated_at       TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT now(), -- Дата последнего обновления платежа (a_money_moves.updated_dt → updated_at (from MySQL))
    accepted_at      TIMESTAMP WITH TIME ZONE,                         -- Дата перехода в один из финальных статусов (a_money_moves.accepted_dt → accepted_at (from MySQL))
    operation_type   PAYMENT_TYPE              NOT NULL,               -- Тип операции платежа (специальный тип для типов)
    code             TEXT                      NOT NULL,               -- Код операции (a_money_moves.payment_system → code (from MySQL))
    account_id       BIGINT                    NOT NULL,               -- Номер счета (a_money_moves.a_id → account_id (from MySQL))
    amount           BIGINT                    NOT NULL,               -- Сумма операции в валюте операции (a_money_moves.amount_in_request_currency → amount (from MySQL))
    amount_currency  TEXT                      NOT NULL,               -- Валюта операции (a_money_moves.to_a_currency → account_currency (from MySQL))
    account_amount   BIGINT                    NOT NULL,               -- Сумма в валюте счета (a_money_moves.amount_in_account_currency)
    exchange_rate    NUMERIC(20,15)            NOT NULL,               -- Обменный курс из валюты операции в валюту счета (a_money_moves.exchange_rate → exchange_rate (from MySQL))
    status           OPERATION_STATUS NOT NULL DEFAULT 'new',          -- Специальный тип для статусов (a_money_moves.status → status (from MySQL))
    comment          TEXT                      NOT NULL,               -- Коментарий необходимый для 1С (a_money_moves.comment → comment (from MySQL))
    mt_order_id      BIGINT,                                           -- ID сделки в Метаке (a_money_moves.mt_order → mt_order_id (from MySQL))
    transaction_id   BIGINT,                                           -- ID транзакции из PS/IN/PIN....
    extra_data       JSONB                     NOT NULL DEFAULT '{}',  -- Данные необходимые для платежа
    purse            TEXT                      NOT NULL DEFAULT '',    -- Кошелек пользователя
    key              TEXT                      NOT NULL DEFAULT ''     -- Спец. поле в том числе и для поиска транзакций
);

CREATE INDEX IF NOT EXISTS idx_payments_operation_type
    ON payments (operation_type);

CREATE INDEX IF NOT EXISTS idx_payments_transaction_id
    ON payments (transaction_id);

CREATE INDEX IF NOT EXISTS idx_payments_account_id
    ON payments (account_id);

-- +goose Down

DROP TABLE IF EXISTS payments;

DROP TYPE IF EXISTS OPERATION_STATUS;
DROP TYPE IF EXISTS PAYMENT_TYPE;
