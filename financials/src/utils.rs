use chrono::{DateTime, Utc, Datelike, TimeZone, Duration};
use chronoutil::delta::shift_months;
use std::error::Error;

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

pub mod yymmdd_format {
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
        let s = format!("{} 08:00:00", s);
        Utc.datetime_from_str(&s, FORMAT).map_err(serde::de::Error::custom)
    }
}

pub fn get_last_day_of_month(dt: DateTime<Utc>) -> DateTime<Utc> {
    let dt = shift_months(dt, 1);
    Utc.with_ymd_and_hms(
        dt.year(),
        dt.month(),
        1,
        8, 0, 0
    ).unwrap() - Duration::days(1)
}

pub fn is_last_day_of_month(dt: DateTime<Utc>) -> bool {
    let next_day = dt + Duration::days(1);
    next_day.day() == 1
}
