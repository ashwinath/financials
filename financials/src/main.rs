use config::Config;
use models::{read_from_csv, Trade, Expense, Asset, Income};
use schema::trades::dsl::trades;
use schema::assets::dsl::assets;
use schema::incomes::dsl::incomes;
use schema::expenses::dsl::expenses;
use stock::calculate_stocks;
use std::error::Error;
use std::process;

#[macro_use]
extern crate diesel;

use chrono::Utc;
use diesel::{Connection, insert_into, delete, RunQueryDsl};
use diesel::pg::PgConnection;
use diesel_migrations::run_pending_migrations;

mod config;
mod models;
mod schema;
mod stock;
mod alphavantage;

fn main() {
    let start_time = Utc::now().time();
    let c = Config::new();

    let conn = init_db(&c.database_url);
    if let Err(e) = load_data(&conn, &c) {
        eprintln!("failed to load csv data: {}", e);
        process::exit(1);
    }

    if let Err(e) = calculate_stocks(&conn, &c.alphavantage_key) {
        eprintln!("failed to load process stocks: {}", e);
        process::exit(1);
    }

    let end_time = Utc::now().time();
    let diff = end_time - start_time;
    println!("Total time taken to run is {} ms.", diff.num_milliseconds());
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
