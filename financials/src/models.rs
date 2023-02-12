use chrono::{DateTime, Utc};
use serde::Deserialize;

use crate::schema::{
    assets,
    average_expenditures,
    exchange_rates,
    expenses,
    incomes,
    shared_expense,
    stocks,
    symbols,
    trades,
    portfolios,
    mortgage,
};

use crate::utils::yymmdd_format;

#[derive(Debug, Deserialize, Queryable, Insertable)]
#[diesel(table_name = trades)]
pub struct TradeWithId {
    pub id: i32,
    pub date_purchased: DateTime<Utc>,
    pub symbol: String,
    pub price_each: f64,
    pub quantity: f64,
    pub trade_type: String,
}

#[derive(Debug, Deserialize, Insertable)]
#[diesel(table_name = trades)]
pub struct Trade {
    #[serde(with = "yymmdd_format")]
    pub date_purchased: DateTime<Utc>,
    pub symbol: String,
    pub price_each: f64,
    pub quantity: f64,
    pub trade_type: String,
}

#[derive(Debug, Deserialize, Queryable, Insertable, Clone)]
#[diesel(table_name = assets)]
pub struct Asset {
    pub id: Option<i32>,
    #[serde(with = "yymmdd_format")]
    #[serde(rename(deserialize = "date"))]
    pub transaction_date: DateTime<Utc>,
    #[serde(rename(deserialize = "type"))]
    pub type_: String,
    pub amount: f64,
}

#[derive(Debug, Deserialize, Queryable, Insertable)]
#[diesel(table_name = incomes)]
pub struct Income {
    pub id: Option<i32>,
    #[serde(with = "yymmdd_format")]
    #[serde(rename(deserialize = "date"))]
    pub transaction_date: DateTime<Utc>,
    #[serde(rename(deserialize = "type"))]
    pub type_: String,
    pub amount: f64,
}

#[derive(Debug, Deserialize, Queryable, Insertable)]
#[diesel(table_name = expenses)]
pub struct Expense {
    pub id: Option<i32>,
    #[serde(with = "yymmdd_format")]
    #[serde(rename(deserialize = "date"))]
    pub transaction_date: DateTime<Utc>,
    #[serde(rename(deserialize = "type"))]
    pub type_: String,
    pub amount: f64,
}

// Need to create another struct for inserting without id
// https://github.com/diesel-rs/diesel/issues/1440
#[derive(Debug, Queryable, Insertable)]
#[diesel(table_name = symbols)]
pub struct SymbolWithId {
    pub id: i32,
    pub symbol_type: String,
    pub symbol: String,
    pub base_currency: Option<String>,
    pub last_processed_date: Option<DateTime<Utc>>,
}

#[derive(Debug, Insertable)]
#[diesel(table_name = symbols)]
pub struct Symbol {
    pub symbol_type: String,
    pub symbol: String,
    pub base_currency: Option<String>,
    pub last_processed_date: Option<DateTime<Utc>>,
}

#[derive(Debug, Insertable, Queryable)]
#[diesel(table_name = exchange_rates)]
pub struct ExchangeRate {
    pub trade_date: DateTime<Utc>,
    pub symbol: String,
    pub price: f64,
}

#[derive(Debug, Insertable)]
#[diesel(table_name = stocks)]
pub struct Stock {
    pub trade_date: DateTime<Utc>,
    pub symbol: String,
    pub price: f64,
}

#[derive(Debug, Insertable, Queryable)]
#[diesel(table_name = portfolios)]
pub struct Portfolio {
    pub trade_date: DateTime<Utc>,
    pub symbol: String,
    pub principal: f64,
    pub nav: f64,
    pub simple_returns: f64,
    pub quantity: f64,
}

#[derive(Debug, Insertable)]
#[diesel(table_name = average_expenditures)]
pub struct AverageExpenditure {
    pub id: Option<i32>,
    pub expense_date: DateTime<Utc>,
    pub amount: f64,
}

#[derive(Debug, PartialEq, Deserialize, Insertable, Clone)]
#[diesel(table_name = mortgage)]
pub struct MortgageSchedule {
    #[serde(with = "yymmdd_format")]
    pub date: DateTime<Utc>,
    pub interest_paid: f64,
    pub principal_paid: f64,
    pub total_principal_paid: f64,
    pub total_interest_paid: f64,
    pub total_principal_left: f64,
    pub total_interest_left: f64,
}

#[derive(Debug, PartialEq, Deserialize, Queryable, Insertable)]
#[diesel(table_name = mortgage)]
pub struct MortgageScheduleWithId {
    pub id: i32,
    #[serde(with = "yymmdd_format")]
    pub date: DateTime<Utc>,
    pub interest_paid: f64,
    pub principal_paid: f64,
    pub total_principal_paid: f64,
    pub total_interest_paid: f64,
    pub total_principal_left: f64,
    pub total_interest_left: f64,
}

#[derive(Debug, Deserialize, Insertable)]
#[diesel(table_name = shared_expense)]
pub struct SharedExpense {
    #[serde(with = "yymmdd_format")]
    #[serde(rename(deserialize = "date"))]
    pub expense_date: DateTime<Utc>,
    #[serde(rename(deserialize = "type"))]
    pub type_: String,
    pub amount: f64,
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::{Utc, TimeZone};
    use crate::utils::read_from_csv;

    #[test]
    fn parse_trades_csv() {
        let result: Vec<Trade> = read_from_csv("./sample/trades.csv").unwrap();
        assert_eq!(result.len(), 12);

        let result = &result[0];
        assert_eq!(result.symbol, "CSPX.LON");
        assert_eq!(result.trade_type, "buy");
        assert_eq!(result.price_each, 446.12);
        assert_eq!(result.quantity, 2.0);
        let expected_date: DateTime<Utc> = Utc
            .datetime_from_str("2021-08-19 08:00:00", "%Y-%m-%d %H:%M:%S")
            .unwrap();
        assert_eq!(result.date_purchased, expected_date);
    }

    #[test]
    fn parse_assets_csv() {
        let result: Vec<Asset> = read_from_csv("./sample/assets.csv").unwrap();
        assert_eq!(result.len(), 76);

        // Asserting the first row should be good enough
        let result = &result[0];
        assert_eq!(result.type_, "Bank");
        assert_eq!(result.amount, 10000.0);
        let expected_date: DateTime<Utc> = Utc
            .datetime_from_str("2021-08-01 08:00:00", "%Y-%m-%d %H:%M:%S")
            .unwrap();
        assert_eq!(result.transaction_date, expected_date);
    }

    #[test]
    fn parse_incomes_csv() {
        let result: Vec<Income> = read_from_csv("./sample/income.csv").unwrap();
        assert_eq!(result.len(), 40);

        // Asserting the first row should be good enough
        let result = &result[0];
        assert_eq!(result.type_, "Base");
        assert_eq!(result.amount, 5000.0);
        let expected_date: DateTime<Utc> = Utc
            .datetime_from_str("2021-08-05 08:00:00", "%Y-%m-%d %H:%M:%S")
            .unwrap();
        assert_eq!(result.transaction_date, expected_date);
    }

    #[test]
    fn parse_expenses_csv() {
        let result: Vec<Expense> = read_from_csv("./sample/expenses.csv").unwrap();
        assert_eq!(result.len(), 77);

        // Asserting the first row should be good enough
        let result = &result[0];
        assert_eq!(result.type_, "Credit Card");
        assert_eq!(result.amount, 1000.0);
        let expected_date: DateTime<Utc> = Utc
            .datetime_from_str("2021-08-31 08:00:00", "%Y-%m-%d %H:%M:%S")
            .unwrap();
        assert_eq!(result.transaction_date, expected_date);
    }

    #[test]
    fn parse_shared_expense_csv() {
        let result: Vec<SharedExpense> = read_from_csv("./sample/shared_expenses.csv").unwrap();
        assert_eq!(result.len(), 40);

        let result = &result[0];
        assert_eq!(result.type_, "Electricity");
        assert_eq!(result.amount, 35.0);
        let expected_date: DateTime<Utc> = Utc
            .datetime_from_str("2022-08-31 08:00:00", "%Y-%m-%d %H:%M:%S")
            .unwrap();
        assert_eq!(result.expense_date, expected_date);
    }
}
