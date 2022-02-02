DROP TABLE IF EXISTS trades;

DROP TABLE IF EXISTS incomes;
DROP INDEX IF EXISTS ix_date_type_incomes;

DROP TABLE IF EXISTS assets;
DROP INDEX IF EXISTS ix_date_type_assets;

DROP TABLE IF EXISTS expenses;
DROP INDEX IF EXISTS ix_date_type_expenses;

DROP TABLE IF EXISTS symbols;

DROP TABLE IF EXISTS exchange_rates;
DROP INDEX IF EXISTS uidx_exchange_rates;

DROP TABLE IF EXISTS stocks;
DROP INDEX IF EXISTS uidx_stocks;
