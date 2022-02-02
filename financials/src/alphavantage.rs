use std::fmt;
use std::error::Error;
use serde::Deserialize;
use reqwest;

//const FX_URL_FORMAT: &str = "https://www.alphavantage.co/query?function=FX_DAILY&from_symbol={}&to_symbol=SGD&apikey={}&outputsize={}";
//const STOCK_URL_FORMAT: &str = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol={}&apikey=%s&outputsize={}";

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

pub fn search_alphavantage_symbol(symbol: &str, api_key: &str) -> Result<AlphaVantageBestMatches, Box<dyn Error>> {
    let url = format!(
        "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords={}&apikey={}",
        symbol,
        api_key,
    );

    let response = reqwest::blocking::get(url)?;

    match response.status() {
        reqwest::StatusCode::OK => return Ok(response.json::<AlphaVantageBestMatches>()?),
        e => return Err(AlphaVantageError {message: format!("status code: {} {}", e.as_str(), e.canonical_reason().unwrap())}.into()),
    }
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
}
