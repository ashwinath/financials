CREATE TABLE IF NOT EXISTS trades (
    id             text NOT NULL PRIMARY KEY,
    user_id        text NOT NULL,
    date_purchased timestamptz NOT NULL,
    symbol         text NOT NULL,
    price_each     double precision NOT NULL,
    quantity       double precision NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trades_id_user_id ON trades(id, user_id);
