CREATE TABLE IF NOT EXISTS trades (
    id             serial NOT NULL PRIMARY KEY,
    date_purchased timestamptz NOT NULL,
    symbol         text NOT NULL,
    price_each     double precision NOT NULL,
    quantity       double precision NOT NULL,
    trade_type     text NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
