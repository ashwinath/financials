use clap::Parser;

#[derive(Parser, Debug)]
#[clap(version)]
pub struct Config {
    #[clap(long)]
    pub database_url: String,
    #[clap(long)]
    pub stocks_csv: String,
}

impl Config {
    pub fn new() -> Config {
        Config::parse()
    }
}
