use std::error::Error;
use crate::alphavantage::{search_alphavantage_symbol, get_currency_history, get_stock_history};
use crate::models::{TradeWithId, Symbol, SymbolWithId, ExchangeRate, Stock, Portfolio};
use crate::schema::trades::dsl::trades;
use crate::schema::symbols::dsl::symbols;
use crate::schema::exchange_rates::dsl::exchange_rates;
use crate::schema::portfolios::dsl::portfolios;
use crate::schema::stocks::dsl::stocks;
use std::collections::HashMap;

use chrono::{DateTime, Duration, Datelike, Utc, TimeZone};
use diesel::{ExpressionMethods, QueryDsl, RunQueryDsl};
use diesel::{insert_into, update};
use diesel::pg::PgConnection;

const FORMAT: &str = "%Y-%m-%d %H:%M:%S";
const STOCK_SYMBOL: &str = "stock";
const CURRENCY_SYMBOL: &str = "currency";

pub fn calculate_stocks(conn: &PgConnection, alphavantage_key: &str) -> Result<(), Box<dyn Error>>  {
    sync_symbol_table(conn, alphavantage_key)?;
    // TODO: Stocks and currencies can be concurrent
    process_currencies(conn, alphavantage_key)?;
    process_stocks(conn, alphavantage_key)?;
    calculate_portfolio(conn)?;

    Ok(())
}

fn calculate_portfolio(conn: &PgConnection) -> Result<(), Box<dyn Error>> {
    let stock_symbols = symbols
        .filter(crate::schema::symbols::dsl::symbol_type.eq(STOCK_SYMBOL))
        .load::<SymbolWithId>(conn)?;

    for stock_symbol in stock_symbols {
        let symbol = stock_symbol.symbol.clone();
        let ts = trades
            .order(crate::schema::trades::date_purchased.asc())
            .filter(crate::schema::trades::symbol.eq(&symbol))
            .load::<TradeWithId>(conn)?;


        let mut partial_portfolios: Vec<Portfolio> = Vec::new();
        let currency_symbol = stock_symbol.base_currency.unwrap();

        // First pass to fill active trading parts
        for t in ts {
            let exchange_rate = get_currency_rate(conn, t.date_purchased, &currency_symbol);
            let trade_multiplier = if t.trade_type == "buy" { 1.0 } else { -1.0 };
            let portfolio = if partial_portfolios.is_empty() {
                let principal = t.price_each * t.quantity * exchange_rate;
                Portfolio {
                    trade_date: t.date_purchased,
                    symbol: symbol.clone(),
                    principal,
                    nav: 0.0, // Calculate later
                    simple_returns: 0.0, // Calculate later
                    quantity: t.quantity,
                }
            } else {
                let last_portfolio = &partial_portfolios[partial_portfolios.len() - 1];
                let principal = last_portfolio.principal + (t.price_each * t.quantity * exchange_rate * trade_multiplier);
                Portfolio {
                    trade_date: t.date_purchased,
                    symbol: symbol.clone(),
                    principal,
                    nav: 0.0, // Calculate later
                    simple_returns: 0.0, // Calculate later
                    quantity: last_portfolio.quantity + (t.quantity * trade_multiplier),
                }
            };
            partial_portfolios.push(portfolio);
        }

        // Second pass to update all gaps in non active trading days
        let mut current_date = partial_portfolios[0].trade_date.clone();
        let mut portfolio_map: HashMap<DateTime<Utc>, Portfolio> = HashMap::new();
        for mut p in partial_portfolios {
            if let Some(p_in_same_day) = portfolio_map.get(&p.trade_date) {
                // There might be multiple trades in a single day for each symbol, we need to combine them
                p.quantity = p_in_same_day.quantity;
                p.principal = p_in_same_day.principal;
            }
            portfolio_map.insert(p.trade_date, p);
        }

        let mut all_portfolios: Vec<Portfolio> = Vec::new();
        let today = chrono::offset::Utc::now();
        let tomorrow = Utc.ymd(today.year(), today.month(), today.day()).and_hms(16, 0, 0);
        while current_date < tomorrow {
            let exchange_rate = get_currency_rate(conn, current_date, &currency_symbol);
            let price = get_stock_price(conn, current_date, &symbol);

            let previous_portfolio = if let Some(previous_portfolio) = portfolio_map.get(&current_date) {
                previous_portfolio
            } else {
                // Guaranteed to have an element.
                all_portfolios.last().unwrap()
            };

            let principal = previous_portfolio.principal;
            let quantity = previous_portfolio.quantity;
            let nav = quantity * price * exchange_rate;
            let simple_returns = (nav - principal) / principal;

            let new_portfolio = Portfolio {
                trade_date: current_date,
                symbol: symbol.clone(),
                principal,
                nav,
                simple_returns,
                quantity,
            };

            current_date = current_date + Duration::days(1);
            all_portfolios.push(new_portfolio);
        }

        insert_into(portfolios)
            .values(&all_portfolios)
            .on_conflict_do_nothing()
            .execute(conn)?;
    }

    Ok(())
}

fn get_stock_price(conn: &PgConnection, trade_date: DateTime<Utc>, symbol: &str) -> f64 {
    let mut trade_date = trade_date;
    loop {
        let price = stocks
            .select(crate::schema::stocks::price)
            .filter(crate::schema::stocks::symbol.eq(symbol))
            .filter(crate::schema::stocks::trade_date.eq(trade_date))
            .first::<f64>(conn);

        if let Ok(value) = price {
            return value;
        }

        trade_date = trade_date - Duration::days(1);
    }
}

fn get_currency_rate(conn: &PgConnection, trade_date: DateTime<Utc>, symbol: &str) -> f64 {
    let mut trade_date = trade_date;
    loop {
        let exchange_rate = exchange_rates
            .select(crate::schema::exchange_rates::price)
            .filter(crate::schema::exchange_rates::symbol.eq(symbol))
            .filter(crate::schema::exchange_rates::trade_date.eq(trade_date))
            .first::<f64>(conn);
        if let Ok(value) = exchange_rate {
            return value;
        }

        trade_date = trade_date - Duration::days(1);
    }
}

fn process_stocks(conn: &PgConnection, alphavantage_key: &str) -> Result<(), Box<dyn Error>> {
    let stock_symbols = symbols
        .filter(crate::schema::symbols::dsl::symbol_type.eq(STOCK_SYMBOL))
        .load::<SymbolWithId>(conn)?;

    // TODO: Can be concurrent
    for stock_symbol in stock_symbols {
        let is_compact = stock_symbol.last_processed_date.is_some();
        let stock_history = get_stock_history(&stock_symbol.symbol, is_compact, alphavantage_key)?
            .results;

        let mut last_processed_date: Option<DateTime<Utc>> = None;
        let mut stock_histories = Vec::new();

        for (date, value) in stock_history {
            let date = format!("{} 16:00:00", date);
            let date = Utc.datetime_from_str(&date, FORMAT)?;
            last_processed_date = if last_processed_date.is_none() {
                Some(date)
            } else if last_processed_date.unwrap().lt(&date) {
                Some(date)
            } else {
                last_processed_date
            };

            let value = value.close;
            let er = Stock {
                trade_date: date,
                symbol: stock_symbol.symbol.to_string(),
                price: value,
            };
            stock_histories.push(er);
        }

        insert_into(stocks)
            .values(&stock_histories)
            .on_conflict_do_nothing()
            .execute(conn)?;

        update(symbols.filter(crate::schema::symbols::dsl::id.eq(stock_symbol.id)))
            .set(crate::schema::symbols::dsl::last_processed_date.eq(last_processed_date))
            .execute(conn)?;
    }

    Ok(())
}

fn process_currencies(conn: &PgConnection, alphavantage_key: &str) -> Result<(), Box<dyn Error>> {
    let currencies = symbols
        .filter(crate::schema::symbols::dsl::symbol_type.eq(CURRENCY_SYMBOL))
        .load::<SymbolWithId>(conn)?;

    // TODO: Can be concurrent
    for currency in currencies {
        let is_compact = currency.last_processed_date.is_some();
        let history = get_currency_history(&currency.symbol, "SGD", is_compact, alphavantage_key)?
            .results;

        let mut last_processed_date: Option<DateTime<Utc>> = None;
        let mut currency_history = Vec::new();

        for (date, value) in history {
            let date = format!("{} 16:00:00", date);
            let date = Utc.datetime_from_str(&date, FORMAT)?;
            last_processed_date = if last_processed_date.is_none() {
                Some(date)
            } else if last_processed_date.unwrap().lt(&date) {
                Some(date)
            } else {
                last_processed_date
            };

            let value = value.close;
            let er = ExchangeRate {
                trade_date: date,
                symbol: currency.symbol.to_string(),
                price: value,
            };
            currency_history.push(er);
        }

        insert_into(exchange_rates)
            .values(&currency_history)
            .on_conflict_do_nothing()
            .execute(conn)?;

        update(symbols.filter(crate::schema::symbols::dsl::id.eq(currency.id)))
            .set(crate::schema::symbols::dsl::last_processed_date.eq(last_processed_date))
            .execute(conn)?;
    }

    Ok(())
}

fn sync_symbol_table(conn: &PgConnection, alphavantage_key: &str) -> Result<(), Box<dyn Error>> {
    let s = trades
        .select(crate::schema::trades::dsl::symbol)
        .distinct()
        .load::<String>(conn)?;

    for symbol in s {
        // Get all missing stock symbols
        let count = symbols
            .filter(crate::schema::symbols::dsl::symbol_type.eq(STOCK_SYMBOL))
            .filter(crate::schema::symbols::dsl::symbol.eq(&symbol))
            .count()
            .get_result::<i64>(conn)?;

        if count == 0 {
            let symbol_info = search_alphavantage_symbol(&symbol, alphavantage_key)?;
            let base_currency = &symbol_info
                .best_matches[0]
                .currency;

            let symbol_object = Symbol {
                symbol_type: STOCK_SYMBOL.to_string(),
                symbol: symbol.to_string(),
                base_currency: Some(base_currency.to_string()),
                last_processed_date: None,
            };

            insert_into(symbols)
                .values(&symbol_object)
                .execute(conn)?;
        }
    }

    // Get all missing currency symbols
    let currencies = symbols
        .select(crate::schema::symbols::dsl::base_currency)
        .filter(crate::schema::symbols::dsl::symbol_type.eq(STOCK_SYMBOL))
        .distinct()
        .load::<Option<String>>(conn)?;

    for currency in currencies {
        // Get all missing stock symbols
        let currency = currency.unwrap(); // Guaranteed to be here

        let count = symbols
            .filter(crate::schema::symbols::dsl::symbol_type.eq(CURRENCY_SYMBOL))
            .filter(crate::schema::symbols::dsl::symbol.eq(&currency))
            .count()
            .get_result::<i64>(conn)?;

        if count == 0 {
            let symbol_object = Symbol {
                symbol_type: CURRENCY_SYMBOL.to_string(),
                symbol: currency.to_string(),
                base_currency: None,
                last_processed_date: None,
            };
            insert_into(symbols)
                .values(&symbol_object)
                .execute(conn)?;
        }
    }

    Ok(())
}
