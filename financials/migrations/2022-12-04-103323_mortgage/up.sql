CREATE TABLE IF NOT EXISTS mortgage (
    id                   serial NOT NULL PRIMARY KEY,
    date                 timestamptz NOT NULL,
    interest_paid        double precision NOT NULL,
    principal_paid       double precision NOT NULL,
    total_principal_paid double precision NOT NULL,
    total_interest_paid  double precision NOT NULL,
    total_principal_left double precision NOT NULL,
    total_interest_left  double precision NOT NULL
);

CREATE UNIQUE INDEX ix_mortgage_date ON mortgage(date);
