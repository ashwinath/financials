CREATE TABLE IF NOT EXISTS assets (
    id               text NOT NULL PRIMARY KEY,
    transaction_date timestamptz NOT NULL,
    type             text NOT NULL,
    amount           double precision not null,
    created_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX ix_date_type_assets ON assets(transaction_date, type);
