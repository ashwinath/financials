use clap::Parser;

#[derive(Parser, Debug)]
#[clap(version)]
pub struct Config {
    #[clap(long)]
    pub stocks_csv: String,
    #[clap(long)]
    pub database_url: String,
}

impl Config {
    pub fn new() -> Config {
        Config::parse()
    }
}
