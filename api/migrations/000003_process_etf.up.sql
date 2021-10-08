CREATE TABLE IF NOT EXISTS symbols (
    id                  text NOT NULL PRIMARY KEY,
    symbol_type         text NOT NULL,
    symbol              text NOT NULL,
    base_currency       text,
    last_processed_date timestamptz,
    created_at          timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_symbols_symbol ON symbols(symbol);

CREATE TABLE IF NOT EXISTS stocks (
    id         text NOT NULL PRIMARY KEY,
    trade_date timestamptz NOT NULL,
    symbol     text NOT NULL,
    price      double precision not null,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS exchange_rates (
    id         text NOT NULL PRIMARY KEY,
    trade_date timestamptz NOT NULL,
    symbol     text NOT NULL,
    price      double precision not null,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);