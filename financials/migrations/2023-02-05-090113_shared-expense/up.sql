CREATE TABLE IF NOT EXISTS shared_expense (
    id             serial NOT NULL PRIMARY KEY,
    expense_date   timestamptz NOT NULL,
    type           text NOT NULL,
    amount         double precision NOT NULL
);

CREATE UNIQUE INDEX ix_shared_expense_type_expense_date ON shared_expense(expense_date, type);
