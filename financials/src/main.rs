use csv;
use config::Config;
use input::Trade;

mod config;
mod input;

fn main() {
    let c = Config::new();
    println!("{}", c.stocks_csv);

    let mut rdr = csv::Reader::from_path(c.stocks_csv).unwrap();

    for result in rdr.deserialize() {
        let record: Trade = result.unwrap();
        println!("{:?}", record)
    }
}
