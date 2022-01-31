use chrono::{DateTime, Utc};
use serde::Deserialize;

use std::error::Error;
use crate::schema::trades;

#[derive(Debug, Deserialize, Queryable, Insertable)]
#[table_name = "trades"]
pub struct Trade {
    pub id: Option<i32>,
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

pub fn read_from_csv<T>(csv: &str) -> Result<Vec<T>, Box<dyn Error>>
    where
    T: for<'de> serde::Deserialize<'de>,
{

    let mut rdr = csv::Reader::from_path(csv)?;

    let mut data = Vec::new();

    for result in rdr.deserialize() {
        let record: T = result?;
        data.push(record);
    }

    Ok(data)
}


#[cfg(test)]
mod tests {
    use super::*;
    use chrono::{Utc, TimeZone};

    #[test]
    fn parse_csv() {
        let result: Vec<Trade> = read_from_csv("./test/trade.csv").unwrap();
        assert_eq!(result.len(), 1);

        let result = &result[0];
        assert_eq!(result.symbol, "IWDA.LON");
        assert_eq!(result.trade_type, "buy");
        assert_eq!(result.price_each, 76.34);
        assert_eq!(result.quantity, 10.0);
        let expected_date: DateTime<Utc> = Utc
            .datetime_from_str("2021-03-11 16:00:00", "%Y-%m-%d %H:%M:%S")
            .unwrap();
        assert_eq!(result.date_purchased, expected_date);
    }
}
