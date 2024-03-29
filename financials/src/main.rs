use asset::{populate_housing_value, populate_investments};
use config::Config;
use expenditure::calculate_average_expenditure;
use models::{Asset, Expense, Income, SharedExpense, Trade};
use mortgage::generate_mortgage_schedule;
use schema::assets::dsl::assets;
use schema::expenses::dsl::expenses;
use schema::incomes::dsl::incomes;
use schema::shared_expense::dsl::shared_expense;
use schema::trades::dsl::trades;
use shared_expenditure::populate_shared_expenditure;
use std::error::Error;
use std::process;
use stock::calculate_stocks;
use telegram_data::load_telegram_data;
use utils::read_from_csv;

#[macro_use]
extern crate diesel;

use chrono::Utc;
use diesel::pg::PgConnection;
use diesel::{delete, insert_into, Connection, RunQueryDsl};
use diesel_migrations::{embed_migrations, EmbeddedMigrations, MigrationHarness};

mod alphavantage;
mod asset;
mod config;
mod expenditure;
mod models;
mod mortgage;
mod schema;
mod shared_expenditure;
mod stock;
mod telegram_data;
mod utils;

const MIGRATIONS: EmbeddedMigrations = embed_migrations!();

fn main() {
    let start_time = Utc::now().time();
    let c = Config::new();

    let mut conn = PgConnection::establish(&c.database_url)
        .expect(&format!("Error connecting to {}", &c.database_url));
    conn.run_pending_migrations(MIGRATIONS).unwrap();

    if let Err(e) = load_data(&mut conn, &c) {
        eprintln!("failed to load csv data: {}", e);
        process::exit(1);
    }

    if let Err(e) = load_telegram_data(&mut conn, &c.telegram_bot_endpoint) {
        eprintln!("failed to download data from telegram bot: {}", e);
        process::exit(8);
    }

    if let Err(e) = calculate_stocks(&mut conn, &c.alphavantage_key) {
        eprintln!("failed to load process stocks: {}", e);
        process::exit(2);
    }

    if let Err(e) = populate_investments(&mut conn) {
        eprintln!("failed to load populate investments into assets: {}", e);
        process::exit(3);
    }

    if let Err(e) = populate_shared_expenditure(&mut conn) {
        eprintln!("failed to load populate shared expenditure: {}", e);
        process::exit(4);
    }

    if let Err(e) = calculate_average_expenditure(&mut conn) {
        eprintln!("failed to load calculate average expenditure: {}", e);
        process::exit(5);
    }

    if let Err(e) = generate_mortgage_schedule(&mut conn, &c.mortgage_yaml) {
        eprintln!("failed to generate mortgage schedule: {}", e);
        process::exit(6);
    }

    if let Err(e) = populate_housing_value(&mut conn) {
        eprintln!("failed to generate mortgage schedule: {}", e);
        process::exit(7);
    }

    let end_time = Utc::now().time();
    let diff = end_time - start_time;
    println!("Total time taken to run is {} ms.", diff.num_milliseconds());
}

fn load_data(conn: &mut PgConnection, c: &Config) -> Result<(), Box<dyn Error>> {
    delete(assets).execute(conn)?;
    let t: Vec<Asset> = read_from_csv(&c.assets_csv)?;
    insert_into(assets).values(&t).execute(conn)?;

    delete(expenses).execute(conn)?;
    let t: Vec<Expense> = read_from_csv(&c.expenses_csv)?;
    insert_into(expenses).values(&t).execute(conn)?;

    delete(incomes).execute(conn)?;
    let t: Vec<Income> = read_from_csv(&c.income_csv)?;
    insert_into(incomes).values(&t).execute(conn)?;

    delete(trades).execute(conn)?;
    let t: Vec<Trade> = read_from_csv(&c.trades_csv)?;
    insert_into(trades).values(&t).execute(conn)?;

    delete(shared_expense).execute(conn)?;
    let t: Vec<SharedExpense> = read_from_csv(&c.shared_expense_csv)?;
    insert_into(shared_expense).values(&t).execute(conn)?;

    Ok(())
}
