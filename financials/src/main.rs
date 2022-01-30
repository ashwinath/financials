use config::Config;

mod config;

fn main() {
    let c = Config::new();
    println!("{}", c.stocks_csv);
}
