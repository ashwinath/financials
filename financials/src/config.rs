use clap::Parser;

#[derive(Parser, Debug)]
#[clap(version)]
pub struct Config {
    #[clap(long)]
    pub database_url: String,
    #[clap(long)]
    pub alphavantage_key: String,
    #[clap(long)]
    pub assets_csv: String,
    #[clap(long)]
    pub expenses_csv: String,
    #[clap(long)]
    pub income_csv: String,
    #[clap(long)]
    pub trades_csv: String,
    #[clap(long)]
    pub mortgage_yaml: String,
    #[clap(long)]
    pub shared_expense_csv: String,
}

impl Config {
    pub fn new() -> Config {
        Config::parse()
    }
}
