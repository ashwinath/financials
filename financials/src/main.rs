use config::Config;
use models::{read_from_csv, Trade};

#[macro_use]
extern crate diesel;

use diesel::{Connection, insert_into};
use diesel::pg::PgConnection;
use crate::diesel::RunQueryDsl;
use schema::trades::dsl::trades;
use diesel_migrations::run_pending_migrations;

mod config;
mod models;
mod schema;

fn main() {
    let c = Config::new();

    // TODO: Test CSV, to be removed.
    let t: Vec<Trade> = read_from_csv(&c.stocks_csv).unwrap();
    println!("{:?}", t);

    // Initialise DB
    let conn = init_db(&c.database_url);
    insert_into(trades)
        .values(&t)
        .execute(&conn)
        .unwrap();
}

fn init_db(database_url: &str) -> PgConnection {
    let conn = PgConnection::establish(database_url)
        .expect(&format!("Error connecting to {}", database_url));

    run_pending_migrations(&conn).expect("Error migration database");

    conn
}
