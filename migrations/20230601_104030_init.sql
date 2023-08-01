-- +goose Up
CREATE TYPE status_type AS ENUM ('active', 'archived', 'deleted');
CREATE TABLE IF NOT EXISTS public.accounts
(
    id                 BIGINT                    NOT NULL,
    user_id            BIGINT                    NOT NULL,
    login              BIGINT                    NOT NULL,
    server_id          SMALLINT                  NOT NULL,
    partner_login      BIGINT                    NOT NULL,
    balance_multiplier INT                       NOT NULL,
    currency           TEXT                      NOT NULL,
    tariff             TEXT                      NOT NULL,
    is_demo            BOOL                      NOT NULL,
    status             status_type               NOT NULL,
    created_at         TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at         TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    balance            BIGINT      DEFAULT 0     NOT NULL,
    platform           TEXT        DEFAULT ''    NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS idx_accounts_login_serverid
    ON accounts (login, server_id);

-- +goose Down
DROP TABLE IF EXISTS public.accounts;
DROP TYPE status_type;
