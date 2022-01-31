use config::Config;
use input::{read_from_csv, Trade};

extern crate diesel;

use diesel::Connection;
use diesel::pg::PgConnection;
use diesel_migrations::run_pending_migrations;

mod config;
mod input;

fn main() {
    let c = Config::new();

    // TODO: Test CSV, to be removed.
    let result: Vec<Trade> = read_from_csv(&c.stocks_csv).unwrap();
    println!("{:?}", result);

    // Initialise DB
    let conn = PgConnection::establish(&c.database_url)
        .expect(&format!("Error connecting to {}", c.database_url));
    run_pending_migrations(&conn).expect("Error migration database");
}
