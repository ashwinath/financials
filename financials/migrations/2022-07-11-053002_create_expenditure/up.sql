CREATE TABLE IF NOT EXISTS average_expenditures (
    id             serial NOT NULL PRIMARY KEY,
    expense_date   timestamptz NOT NULL,
    amount         double precision NOT NULL
);

CREATE UNIQUE INDEX ix_average_expenditures_date ON average_expenditures(expense_date);
