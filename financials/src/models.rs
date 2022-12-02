use chrono::{DateTime, Utc};
use serde::Deserialize;

use std::error::Error;
use crate::schema::{
    assets,
    average_expenditures,
    exchange_rates,
    expenses,
    incomes,
    stocks,
    symbols,
    trades,
    portfolios,
};

mod yymmdd_format {
    use chrono::{DateTime, Utc, TimeZone};
    use serde::{self, Deserialize, Deserializer};

    const FORMAT: &'static str = "%Y-%m-%d %H:%M:%S";

    pub fn deserialize<'de, D>(
        deserializer: D,
    ) -> Result<DateTime<Utc>, D::Error>
    where
        D: Deserializer<'de>,
    {
        let s = String::deserialize(deserializer)?;
        let s = format!("{} 08:00:00", s);
        Utc.datetime_from_str(&s, FORMAT).map_err(serde::de::Error::custom)
    }
}

pub fn read_from_csv<T>(csv: &str) -> Result<Vec<T>, Box<dyn Error>>
    where
    T: for<'de> serde::Deserialize<'de>,
{

    let mut rdr = csv::Reader::from_path(csv)?;

    let mut data = Vec::new();

    for result in rdr.deserialize() {
        let record: T = result?;
        data.push(record);
    }

    Ok(data)
}
#[derive(Debug, Deserialize, Queryable, Insertable)]
#[table_name = "trades"]
pub struct TradeWithId {
    pub id: i32,
    pub date_purchased: DateTime<Utc>,
    pub symbol: String,
    pub price_each: f64,
    pub quantity: f64,
    pub trade_type: String,
}

#[derive(Debug, Deserialize, Insertable)]
#[table_name = "trades"]
pub struct Trade {
    #[serde(with = "yymmdd_format")]
    pub date_purchased: DateTime<Utc>,
    pub symbol: String,
    pub price_each: f64,
    pub quantity: f64,
    pub trade_type: String,
}

#[derive(Debug, Deserialize, Queryable, Insertable)]
#[table_name = "assets"]
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
#[table_name = "incomes"]
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
#[table_name = "expenses"]
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
#[table_name = "symbols"]
pub struct SymbolWithId {
    pub id: i32,
    pub symbol_type: String,
    pub symbol: String,
    pub base_currency: Option<String>,
    pub last_processed_date: Option<DateTime<Utc>>,
}

#[derive(Debug, Insertable)]
#[table_name = "symbols"]
pub struct Symbol {
    pub symbol_type: String,
    pub symbol: String,
    pub base_currency: Option<String>,
    pub last_processed_date: Option<DateTime<Utc>>,
}

#[derive(Debug, Insertable, Queryable)]
#[table_name = "exchange_rates"]
pub struct ExchangeRate {
    pub trade_date: DateTime<Utc>,
    pub symbol: String,
    pub price: f64,
}

#[derive(Debug, Insertable)]
#[table_name = "stocks"]
pub struct Stock {
    pub trade_date: DateTime<Utc>,
    pub symbol: String,
    pub price: f64,
}

#[derive(Debug, Insertable, Queryable)]
#[table_name = "portfolios"]
pub struct Portfolio {
    pub trade_date: DateTime<Utc>,
    pub symbol: String,
    pub principal: f64,
    pub nav: f64,
    pub simple_returns: f64,
    pub quantity: f64,
}

#[derive(Debug, Insertable)]
#[table_name = "average_expenditures"]
pub struct AverageExpenditure {
    pub id: Option<i32>,
    pub expense_date: DateTime<Utc>,
    pub amount: f64,
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::{Utc, TimeZone};

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
        assert_eq!(result.len(), 68);

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
        assert_eq!(result.len(), 34);

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
        assert_eq!(result.len(), 65);

        // Asserting the first row should be good enough
        let result = &result[0];
        assert_eq!(result.type_, "Credit Card");
        assert_eq!(result.amount, 1000.0);
        let expected_date: DateTime<Utc> = Utc
            .datetime_from_str("2021-08-31 08:00:00", "%Y-%m-%d %H:%M:%S")
            .unwrap();
        assert_eq!(result.transaction_date, expected_date);
    }
}
