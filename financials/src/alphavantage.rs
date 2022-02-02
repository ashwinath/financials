use std::fmt;
use std::error::Error;
use std::collections::HashMap;
use serde::{Deserialize, Deserializer};
use serde_json::Value;
use reqwest;

#[derive(Debug, Clone)]
pub struct AlphaVantageError {
    pub message: String,
}

impl Error for AlphaVantageError {}

impl fmt::Display for AlphaVantageError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "AlphaVantage Error")
    }
}

#[derive(Debug, Deserialize)]
pub struct AlphaVantageBestMatches {
    #[serde(rename(deserialize = "bestMatches"))]
    pub best_matches: Vec<AlphaVantageSymbolSearchResult>,
}

#[derive(Debug, Deserialize)]
pub struct AlphaVantageSymbolSearchResult {
    #[serde(rename(deserialize = "1. symbol"))]
    pub symbol: String,
    #[serde(rename(deserialize = "8. currency"))]
    pub currency: String,
}

pub fn search_alphavantage_symbol(symbol: &str, api_key: &str) -> Result<AlphaVantageBestMatches, Box<dyn Error>> {
    let url = format!(
        "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords={}&apikey={}",
        symbol,
        api_key,
    );

    call_alphavantage(&url)
}

#[derive(Debug, Deserialize)]
pub struct AlphaVantageCurrencyResult {
    #[serde(rename(deserialize = "Time Series FX (Daily)"))]
    pub results: HashMap<String, AlphaVantageCurrencyDailyResult>,
}

#[derive(Debug, Deserialize)]
pub struct AlphaVantageCurrencyDailyResult {
    #[serde(deserialize_with = "de_float")]
    #[serde(rename(deserialize = "4. close"))]
    pub close: f64,
}

pub fn get_currency_history(from_symbol: &str, to_symbol: &str, is_compact: bool, api_key: &str) -> Result<AlphaVantageCurrencyResult, Box<dyn Error>> {
    let url = format!(
        "https://www.alphavantage.co/query?function=FX_DAILY&from_symbol={}&to_symbol={}&outputsize={}&apikey={}",
        from_symbol,
        to_symbol,
        if is_compact {"compact"} else {"full"},
        api_key,
    );

    call_alphavantage(&url)
}

#[derive(Debug, Deserialize)]
pub struct AlphaVantageStockResult {
    #[serde(rename(deserialize = "Time Series (Daily)"))]
    pub results: HashMap<String, AlphaVantageStockDailyResult>,
}

#[derive(Debug, Deserialize)]
pub struct AlphaVantageStockDailyResult {
    #[serde(deserialize_with = "de_float")]
    #[serde(rename(deserialize = "4. close"))]
    pub close: f64,
}

pub fn get_stock_history(symbol: &str, is_compact: bool, api_key: &str) -> Result<AlphaVantageStockResult, Box<dyn Error>> {
    let url = format!(
        "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol={}&outputsize={}&apikey={}",
        symbol,
        if is_compact {"compact"} else {"full"},
        api_key,
    );

    call_alphavantage(&url)
}

fn call_alphavantage<T>(url: &str) -> Result<T, Box<dyn Error>> 
    where
    T: for<'de> serde::Deserialize<'de>
{
    let response = reqwest::blocking::get(url)?;

    match response.status() {
        reqwest::StatusCode::OK => return Ok(response.json::<T>()?),
        e => return Err(AlphaVantageError {message: format!("status code: {} {}", e.as_str(), e.canonical_reason().unwrap())}.into()),
    }
}

fn de_float<'de, D: Deserializer<'de>>(deserializer: D) -> Result<f64, D::Error> {
    Ok(match Value::deserialize(deserializer)? {
        Value::String(s) => s.parse().map_err(serde::de::Error::custom)?,
        Value::Number(num) => num.as_f64().ok_or(serde::de::Error::custom("Invalid number"))?,
        _ => return Err(serde::de::Error::custom("wrong type"))
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn search_symbol() {
        let result = search_alphavantage_symbol("tesco", "demo").unwrap();
        assert_eq!(result.best_matches.len(), 3);
        let result = &result.best_matches[0];
        assert_eq!(result.symbol, "TSCO.LON");
        assert_eq!(result.currency, "GBX");
    }

    #[test]
    fn currency_history() {
        let result = get_currency_history("EUR", "USD", false, "demo").unwrap();
        let result = result.results.get("2022-02-02").unwrap();
        assert!(result.close > 0.0);
    }

    #[test]
    fn stock_history() {
        let result = get_stock_history("IBM", false, "demo").unwrap();
        let result = result.results.get("2022-02-01").unwrap();
        assert!(result.close > 0.0);
    }
}
