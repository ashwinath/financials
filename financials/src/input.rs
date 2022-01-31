use chrono::{DateTime, Utc};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Trade {
    #[serde(with = "yymmdd_format")]
    pub date_purchased: DateTime<Utc>,
    pub symbol: String,
    pub trade_type: String,
    pub price_each: f64,
    pub quantity: f64,
}

mod yymmdd_format {
    use chrono::{DateTime, Utc, TimeZone};
    use serde::{self, Deserialize, Deserializer};

    const FORMAT: &'static str = "%Y-%m-%d %H:%M:%S";

    pub fn deserialize<'de, D>(
        deserializer: D,
    ) -> Result<DateTime<Utc>, D::Error>
    where
        D: Deserializer<'de>,
    {
        let s = String::deserialize(deserializer)?;
        let s = format!("{} 16:00:00", s);
        Utc.datetime_from_str(&s, FORMAT).map_err(serde::de::Error::custom)
    }
}

pub fn read_from_csv<T>(csv: &str) -> Vec<T>
    where
    T: for<'de> serde::Deserialize<'de>,
{

    let mut rdr = csv::Reader::from_path(csv).unwrap();

    let mut data = Vec::new();

    for result in rdr.deserialize() {
        let record: T = result.unwrap();
        data.push(record);
    }

    data
}


#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parse_csv() {
        let result: Vec<Trade> = read_from_csv("./test/trade.csv");
        println!("{:?}", result);
    }
}
