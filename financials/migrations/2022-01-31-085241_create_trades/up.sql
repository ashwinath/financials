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

CREATE TABLE IF NOT EXISTS assets (
    id               serial NOT NULL PRIMARY KEY,
    transaction_date timestamptz NOT NULL,
    type             text NOT NULL,
    amount           double precision not null,
    created_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX ix_date_type_assets ON assets(transaction_date, type);

CREATE TABLE IF NOT EXISTS incomes (
    id               serial NOT NULL PRIMARY KEY,
    transaction_date timestamptz NOT NULL,
    type             text NOT NULL,
    amount           double precision not null,
    created_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX ix_date_type_incomes ON incomes(transaction_date, type);

CREATE TABLE IF NOT EXISTS expenses (
    id               serial NOT NULL PRIMARY KEY,
    transaction_date timestamptz NOT NULL,
    type             text NOT NULL,
    amount           double precision not null,
    created_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX ix_date_type_expenses ON expenses(transaction_date, type);
