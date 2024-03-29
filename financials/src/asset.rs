use chrono::{DateTime, Utc, Datelike, TimeZone, Duration};
use chronoutil::delta::shift_months;
use crate::schema::assets::dsl::assets;
use crate::schema::portfolios::dsl::portfolios;
use crate::schema::mortgage::dsl::mortgage;
use std::error::Error;
use crate::models::{Asset, MortgageScheduleWithId};
use diesel::pg::PgConnection;
use diesel::dsl::sum;
use diesel::{insert_into, ExpressionMethods, QueryDsl, RunQueryDsl};

// Populates the investments on a monthly basis, every first of the month,
// check the value of the investments and populate into the assets table.
// Removing the need to manually key in the assets CSV.
pub fn populate_investments(conn: &mut PgConnection) -> Result<(), Box<dyn Error>> {
    // Get earliest start date
    let first_date = portfolios
        .select(crate::schema::portfolios::trade_date)
        .order(crate::schema::portfolios::trade_date.asc())
        .first::<DateTime<Utc>>(conn)?;

    let mut current_date = if first_date.day() == 1 {
        first_date.clone()
    } else {
        shift_months(first_date, 1)
    };

    current_date = Utc.with_ymd_and_hms(
        current_date.year(),
        current_date.month(),
        1,
        8, 0, 0
    ).unwrap();

    let mut all_investments: Vec<Asset> = Vec::new();
    let tomorrow = chrono::offset::Utc::now() + Duration::days(1);
    while current_date < tomorrow {
        let amount = portfolios
            .select(sum(crate::schema::portfolios::nav))
            .filter(crate::schema::portfolios::dsl::trade_date.eq(current_date))
            .first::<Option<f64>>(conn)?;

        if let Some(i) = amount {
            let asset = Asset {
                id: None,
                transaction_date: current_date.clone(),
                type_: String::from("Investments"),
                amount: i,
            };
            all_investments.push(asset);
        }

        current_date = shift_months(current_date, 1);
    }

    insert_into(assets)
        .values(&all_investments)
        .execute(conn)?;

    Ok(())
}

const HOUSE_SPLIT_RATIO: f64 = 2.0;

// Populates the assets of the principal paid in the mortgage
pub fn populate_housing_value(conn: &mut PgConnection) -> Result<(), Box<dyn Error>> {
    let mortgages = mortgage
        .order_by(crate::schema::mortgage::dsl::date.asc())
        .load::<MortgageScheduleWithId>(conn)?;
    let house_assets: Vec<Asset> = mortgages.iter().map(|m| {
        let date = Utc.with_ymd_and_hms(
            m.date.year(),
            m.date.month(),
            1,
            8, 0, 0
        ).unwrap();
        Asset {
            id: None,
            transaction_date: date,
            type_: String::from("House"),
            amount: m.total_principal_paid / HOUSE_SPLIT_RATIO,
        }
    }).collect();

    // Find gaps in between dates that have no mortgage schedule
    if house_assets.len() == 0 {
        return Ok(());
    }

    let mut all_house_assets: Vec<Asset> = Vec::new();
    for counter in 0..house_assets.len() - 1 {
        let house_asset = &house_assets[counter];

        let mut current_date = house_asset.transaction_date;

        let next_house_asset = &house_assets[counter + 1];

        while current_date < next_house_asset.transaction_date {
            let asset = Asset {
                id: None,
                transaction_date: current_date.clone(),
                type_: String::from("House"),
                amount: house_asset.clone().amount,
            };
            current_date = shift_months(current_date, 1);
            all_house_assets.push(asset);
        }
    }

    insert_into(assets)
        .values(&all_house_assets)
        .execute(conn)?;

    Ok(())
}
