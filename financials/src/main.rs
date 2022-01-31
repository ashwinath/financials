use config::Config;
use input::{read_from_csv, Trade};

mod config;
mod input;

fn main() {
    let c = Config::new();

    let result: Vec<Trade> = read_from_csv(&c.stocks_csv).unwrap();
    println!("{:?}", result);
}
