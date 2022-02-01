use config::Config;
use models::{read_from_csv, Trade};
use std::error::Error;
use std::process;

#[macro_use]
extern crate diesel;

use crate::diesel::RunQueryDsl;
use diesel::{Connection, insert_into};
use diesel::pg::PgConnection;
use schema::trades::dsl::trades;
use diesel_migrations::run_pending_migrations;

mod config;
mod models;
mod schema;

fn main() {
    let c = Config::new();

    let conn = init_db(&c.database_url);
    if let Err(e) = load_data(&conn, &c) {
        eprintln!("failed to load csv data: {}", e);
        process::exit(1);
    }
}

fn load_data(conn: &PgConnection, c: &Config) -> Result<(), Box<dyn Error>> {
    // Load Stocks data
    let t: Vec<Trade> = read_from_csv(&c.stocks_csv)?;
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
