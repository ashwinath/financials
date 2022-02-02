CREATE TABLE IF NOT EXISTS trades (
    id             serial NOT NULL PRIMARY KEY,
    date_purchased timestamptz NOT NULL,
    symbol         text NOT NULL,
    price_each     double precision NOT NULL,
    quantity       double precision NOT NULL,
    trade_type     text NOT NULL
);

CREATE TABLE IF NOT EXISTS assets (
    id               serial NOT NULL PRIMARY KEY,
    transaction_date timestamptz NOT NULL,
    type             text NOT NULL,
    amount           double precision not null
);

CREATE UNIQUE INDEX ix_date_type_assets ON assets(transaction_date, type);

CREATE TABLE IF NOT EXISTS incomes (
    id               serial NOT NULL PRIMARY KEY,
    transaction_date timestamptz NOT NULL,
    type             text NOT NULL,
    amount           double precision not null
);

CREATE UNIQUE INDEX ix_date_type_incomes ON incomes(transaction_date, type);

CREATE TABLE IF NOT EXISTS expenses (
    id               serial NOT NULL PRIMARY KEY,
    transaction_date timestamptz NOT NULL,
    type             text NOT NULL,
    amount           double precision not null
);

CREATE UNIQUE INDEX ix_date_type_expenses ON expenses(transaction_date, type);

CREATE TABLE IF NOT EXISTS symbols (
    id                  serial NOT NULL PRIMARY KEY,
    symbol_type         text NOT NULL,
    symbol              text NOT NULL,
    base_currency       text,
    last_processed_date timestamptz
);

CREATE TABLE IF NOT EXISTS exchange_rates (
    id         serial NOT NULL PRIMARY KEY,
    trade_date timestamptz NOT NULL,
    symbol     text NOT NULL,
    price      double precision not null
);

CREATE UNIQUE INDEX uidx_exchange_rates ON exchange_rates(trade_date, symbol);

CREATE TABLE IF NOT EXISTS stocks (
    id         serial NOT NULL PRIMARY KEY,
    trade_date timestamptz NOT NULL,
    symbol     text NOT NULL,
    price      double precision not null
);

CREATE UNIQUE INDEX uidx_stocks ON stocks(trade_date, symbol);

CREATE TABLE IF NOT EXISTS portfolios (
    id             text NOT NULL PRIMARY KEY,
    trade_date     timestamptz NOT NULL,
    symbol         text NOT NULL,
    principal      double precision NOT NULL,
    nav            double precision NOT NULL,
    simple_returns double precision NOT NULL,
    quantity       double precision NOT NULL
);
CREATE UNIQUE INDEX uidx_portfolios ON portfolios(trade_date, symbol);
