use config::Config;
use models::{Trade, Expense, Asset, Income, SharedExpense};
use schema::trades::dsl::trades;
use schema::assets::dsl::assets;
use schema::incomes::dsl::incomes;
use schema::expenses::dsl::expenses;
use schema::shared_expense::dsl::shared_expense;
use stock::calculate_stocks;
use asset::{populate_investments, populate_housing_value};
use expenditure::calculate_average_expenditure;
use mortgage::generate_mortgage_schedule;
use utils::read_from_csv;
use std::error::Error;
use std::process;

#[macro_use]
extern crate diesel;

use chrono::Utc;
use diesel::{Connection, insert_into, delete, RunQueryDsl};
use diesel::pg::PgConnection;
use diesel_migrations::run_pending_migrations;

mod asset;
mod config;
mod expenditure;
mod models;
mod schema;
mod stock;
mod alphavantage;
mod mortgage;
mod utils;

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

    if let Err(e) = populate_investments(&conn) {
        eprintln!("failed to load populate investments into assets: {}", e);
        process::exit(1);
    }

    if let Err(e) = calculate_average_expenditure(&conn) {
        eprintln!("failed to load calculate average expenditure: {}", e);
        process::exit(1);
    }

    if let Err(e) = generate_mortgage_schedule(&conn, &c.mortgage_yaml) {
        eprintln!("failed to generate mortgage schedule: {}", e);
        process::exit(1);
    }

    if let Err(e) = populate_housing_value(&conn) {
        eprintln!("failed to generate mortgage schedule: {}", e);
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

    delete(shared_expense).execute(conn)?;
    let t: Vec<SharedExpense> = read_from_csv(&c.shared_expense_csv)?;
    insert_into(shared_expense)
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
