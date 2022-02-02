use alphavantage::{search_alphavantage_symbol, get_currency_history, get_stock_history};
use config::Config;
use models::{read_from_csv, Trade, Expense, Asset, Income, Symbol, SymbolWithId, ExchangeRate};
use schema::trades::dsl::trades;
use schema::assets::dsl::assets;
use schema::incomes::dsl::incomes;
use schema::expenses::dsl::expenses;
use schema::symbols::dsl::symbols;
use schema::exchange_rates::dsl::exchange_rates;
use std::error::Error;
use std::process;

#[macro_use]
extern crate diesel;

use crate::diesel::{ExpressionMethods, QueryDsl, RunQueryDsl};
use chrono::{Utc, TimeZone};
use diesel::{Connection, insert_into, delete};
use diesel::pg::PgConnection;
use diesel_migrations::run_pending_migrations;

mod alphavantage;
mod config;
mod models;
mod schema;

const FORMAT: &'static str = "%Y-%m-%d %H:%M:%S";

fn main() {
    let start_time = Utc::now().time();
    let c = Config::new();

    let conn = init_db(&c.database_url);
    if let Err(e) = load_data(&conn, &c) {
        eprintln!("failed to load csv data: {}", e);
        process::exit(1);
    }

    if let Err(e) = sync_symbol_table(&conn, &c.alphavantage_key) {
        eprintln!("failed to load csv data: {}", e);
        process::exit(1);
    }

    if let Err(e) = process_currencies(&conn, &c.alphavantage_key) {
        eprintln!("failed to load csv data: {}", e);
        process::exit(1);
    }

    let end_time = Utc::now().time();
    let diff = end_time - start_time;
    println!("Total time taken to run is {} ms.", diff.num_milliseconds());
}

fn process_currencies(conn: &PgConnection, alphavantage_key: &str) -> Result<(), Box<dyn Error>> {
    let currencies = symbols
        .filter(schema::symbols::dsl::symbol_type.eq("currency"))
        .load::<SymbolWithId>(conn)?;

    for currency in currencies {
        let is_compact = currency.last_processed_date.is_some();
        let history = get_currency_history(&currency.symbol, "SGD", is_compact, alphavantage_key)?
            .results;

        let mut currency_history = Vec::new();
        for (date, value) in history {
            let date = format!("{} 16:00:00", date);
            let date = Utc.datetime_from_str(&date, FORMAT)?;
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
    }

    Ok(())
}

fn sync_symbol_table(conn: &PgConnection, alphavantage_key: &str) -> Result<(), Box<dyn Error>> {
    let s = trades
        .select(schema::trades::dsl::symbol)
        .distinct()
        .load::<String>(conn)?;

    for symbol in s {
        // Get all missing stock symbols
        let count = symbols
            .filter(schema::symbols::dsl::symbol_type.eq("stock"))
            .filter(schema::symbols::dsl::symbol.eq(&symbol))
            .count()
            .get_result::<i64>(conn)?;

        if count == 0 {
            let symbol_info = search_alphavantage_symbol(&symbol, alphavantage_key)?;
            let base_currency = &symbol_info
                .best_matches[0]
                .currency;

            let symbol_object = Symbol {
                symbol_type: "stock".to_string(),
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
        .select(schema::symbols::dsl::base_currency)
        .filter(schema::symbols::dsl::symbol_type.eq("stock"))
        .distinct()
        .load::<Option<String>>(conn)?;

    for currency in currencies {
        // Get all missing stock symbols
        let currency = currency.unwrap(); // Guaranteed to be here

        let count = symbols
            .filter(schema::symbols::dsl::symbol_type.eq("currency"))
            .filter(schema::symbols::dsl::symbol.eq(&currency))
            .count()
            .get_result::<i64>(conn)?;

        if count == 0 {
            let symbol_object = Symbol {
                symbol_type: "currency".to_string(),
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

fn load_data(conn: &PgConnection, c: &Config) -> Result<(), Box<dyn Error>> {
    delete(assets).execute(conn)?;
    let t: Vec<Asset> = read_from_csv(&c.assets_csv)?;
    insert_into(assets)
        .values(&t)
        .execute(conn)?;

    delete(expenses).execute(conn)?;
    let t: Vec<Expense> = read_from_csv(&c.expenses_csv)?;
    insert_into(expenses)
        .values(&t)
        .execute(conn)?;

    delete(incomes).execute(conn)?;
    let t: Vec<Income> = read_from_csv(&c.income_csv)?;
    insert_into(incomes)
        .values(&t)
        .execute(conn)?;

    delete(trades).execute(conn)?;
    let t: Vec<Trade> = read_from_csv(&c.trades_csv)?;
    insert_into(trades)
        .values(&t)
        .execute(conn)?;

    Ok(())
}

fn init_db(database_url: &str) -> PgConnection {
    let conn = PgConnection::establish(database_url)
        .expect(&format!("Error connecting to {}", database_url));

    run_pending_migrations(&conn).expect("Error migration database");

    conn
}
