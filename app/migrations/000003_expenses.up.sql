CREATE TABLE IF NOT EXISTS expenses (
    id               text NOT NULL PRIMARY KEY,
    transaction_date timestamptz NOT NULL,
    type             text NOT NULL,
    amount           double precision not null,
    created_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX tx_date_type_expenses ON expenses(transaction_date, type);
