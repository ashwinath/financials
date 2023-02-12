use chrono::{DateTime, Utc};
use chronoutil::delta::shift_months;
use diesel::{insert_into, delete, RunQueryDsl};
use diesel::pg::PgConnection;
use serde::Deserialize;

use crate::utils::yymmdd_format;
use crate::models::MortgageSchedule;
use crate::schema::mortgage::dsl::mortgage;

use std::error::Error;
use std::fs;

const NUMBER_OF_MONTHS_IN_YEAR: i32 = 12;

#[derive(Debug, PartialEq, Deserialize)]
struct MortgageConfig {
    pub mortgages: Vec<Mortgage>,
}

#[derive(Debug, PartialEq, Deserialize)]
struct Mortgage {
    pub total: f64,
    #[serde(with = "yymmdd_format")]
    pub mortgage_first_payment: DateTime<Utc>,
    pub mortgage_duration_in_years: i32,
    pub downpayments: Vec<Downpayment>,
    pub interest_rate_percentage: f64,
}

#[derive(Debug, PartialEq, Deserialize)]
struct Downpayment {
    #[serde(with = "yymmdd_format")]
    pub date: DateTime<Utc>,
    pub sum: f64,
}

pub fn generate_mortgage_schedule(conn: &mut PgConnection, config_location: &str) -> Result<(), Box<dyn Error>> {
    let config = parse_mortgage(config_location)?;
    delete(mortgage).execute(conn)?;
    for m in config.mortgages.iter() {
        let mortgage_schedule = parse_one_mortgage_schedule(m);
        insert_into(mortgage)
            .values(&mortgage_schedule)
            .execute(conn)?;
    }

    Ok(())
}

fn parse_one_mortgage_schedule(m: &Mortgage) -> Vec<MortgageSchedule> {
    let principal = {
        let paid_already: f64 = m.downpayments.iter().map(|x| x.sum).sum();
        m.total - paid_already
    };

    let monthly_payment = calculate_mortgage(
        principal,
        m.interest_rate_percentage,
        m.mortgage_duration_in_years
    );

    let interest_paid_schedule = calculate_interest_paid_schedule(
        principal,
        monthly_payment,
        m.interest_rate_percentage,
    );

    let total_interest_to_be_paid = interest_paid_schedule.iter().sum();

    let mut total_interest_left: f64 = total_interest_to_be_paid;
    let mut total_principal_paid = 0.0;
    let mut total_interest_paid = 0.0;
    let mut total_principal_left = m.total;

    let mut mortgage_schedule: Vec<MortgageSchedule> = Vec::new();

    // Downpayment
    for downpayment in m.downpayments.iter() {
        total_principal_left -= downpayment.sum;
        total_principal_paid += downpayment.sum;

        let schedule = MortgageSchedule {
            date: downpayment.date,
            interest_paid: 0.0,
            principal_paid: downpayment.sum,
            total_principal_paid,
            total_interest_paid,
            total_principal_left,
            total_interest_left,
        };
        mortgage_schedule.push(schedule);
    }

    let mut mortgage_date = m.mortgage_first_payment;

    for interest_paid in interest_paid_schedule.iter() {
        // Interest
        total_interest_paid = {
            let total_interest_paid = total_interest_paid + interest_paid;

            if total_interest_paid > total_interest_to_be_paid {
                total_interest_to_be_paid
            } else {
                total_interest_paid
            }
        };

        total_interest_left = {
            let total_interest_left = total_interest_left - interest_paid;
            if total_interest_left < 0.0 {
                0.0
            } else {
                total_interest_left
            }
        };

        // Principal
        let principal_paid = monthly_payment - interest_paid;
        total_principal_paid = {
            let total_principal_paid = total_principal_paid + principal_paid;
            if total_interest_paid > m.total {
                m.total
            } else {
                total_principal_paid
            }
        };
        total_principal_left = {
            let total_principal_left = total_principal_left - principal_paid;
            if total_principal_left < 0.0 {
                0.0
            } else {
                total_principal_left
            }
        };

        let schedule = MortgageSchedule {
            date: mortgage_date.clone(),
            interest_paid: interest_paid.clone(),
            principal_paid,
            total_principal_paid,
            total_interest_paid,
            total_principal_left,
            total_interest_left,
        };

        mortgage_date = shift_months(mortgage_date, 1);

        mortgage_schedule.push(schedule);
    }

    mortgage_schedule
}

fn parse_mortgage(config_location: &str) -> Result<MortgageConfig, Box<dyn Error>> {
    let yaml_string = fs::read_to_string(config_location)?;
    let config: MortgageConfig = serde_yaml::from_str(&yaml_string)?;
    Ok(config)
}

fn calculate_interest_paid_schedule(principal: f64, monthly_payment: f64, interest_rate: f64) -> Vec<f64> {
    let ir = interest_rate / 100.0 / NUMBER_OF_MONTHS_IN_YEAR as f64;
    let mut sum_left = principal;

    let mut interest_paid_schedule: Vec<f64> = Vec::new();

    while sum_left > 0.0 {
        let interest_paid = sum_left * ir;
        interest_paid_schedule.push(interest_paid);
        sum_left += interest_paid;
        sum_left -= monthly_payment;
    }

    interest_paid_schedule
}

fn calculate_mortgage(principal: f64, interest_rate: f64, mortgage_duration_in_years: i32) -> f64 {
    let number_of_months = mortgage_duration_in_years * NUMBER_OF_MONTHS_IN_YEAR;
    let ir = interest_rate / 100.0 / NUMBER_OF_MONTHS_IN_YEAR as f64;
    // M = P [ i(1 + i)^n ] / [ (1 + i)^n – 1]. 
    let monthly_payment = principal * (
        (ir * f64::powi(1.0 + ir, number_of_months)) /
        (f64::powi(1.0 + ir, number_of_months) - 1.0)
    );
    monthly_payment
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::{Utc, TimeZone};

    #[test]
    fn test_mortgage_schedule() {
        let mortgage_config = parse_mortgage("./sample/mortgage.yaml").unwrap();
        let mortgage_schedule = parse_one_mortgage_schedule(mortgage_config.mortgages.get(0).unwrap());

        let first_downpayment_schedule = mortgage_schedule.get(0).unwrap();
        assert_eq!(
            first_downpayment_schedule.date,
            Utc.datetime_from_str(
                "2021-10-10 08:00:00",
                "%Y-%m-%d %H:%M:%S"
            ).unwrap()
        );
        assert_eq!(format!("{:.2}", first_downpayment_schedule.interest_paid), "0.00");
        assert_eq!(format!("{:.2}", first_downpayment_schedule.total_interest_paid), "0.00");
        assert_eq!(format!("{:.2}", first_downpayment_schedule.principal_paid), "1000.00");
        assert_eq!(format!("{:.2}", first_downpayment_schedule.total_principal_paid), "1000.00");
        assert_eq!(format!("{:.2}", first_downpayment_schedule.total_principal_left), "49000.00");
        assert_eq!(format!("{:.2}", first_downpayment_schedule.total_interest_left), "10469.25");

        let first_downpayment_schedule = mortgage_schedule.get(2).unwrap();
        assert_eq!(
            first_downpayment_schedule.date,
            Utc.datetime_from_str(
                "2022-10-10 08:00:00",
                "%Y-%m-%d %H:%M:%S"
            ).unwrap()
        );
        assert_eq!(format!("{:.2}", first_downpayment_schedule.interest_paid), "62.83");
        assert_eq!(format!("{:.2}", first_downpayment_schedule.total_interest_paid), "62.83");
        assert_eq!(format!("{:.2}", first_downpayment_schedule.principal_paid), "68.73");
        assert_eq!(format!("{:.2}", first_downpayment_schedule.total_principal_paid), "21068.73");
        assert_eq!(format!("{:.2}", first_downpayment_schedule.total_principal_left), "28931.27");
        assert_eq!(format!("{:.2}", first_downpayment_schedule.total_interest_left), "10406.41");
    }

    #[test]
    fn test_mortgage() {
        struct TestCase<'a> {
            pub principal: f64,
            pub ir: f64,
            pub years: i32,
            pub expected: &'a str,
            pub interest_paid_expected: &'a str,
        }

        let test_cases = vec!(
            TestCase {
                principal: 29_000.0,
                ir: 2.6,
                years: 25,
                expected: "131.56",
                interest_paid_expected: "10469.25",
            },
            TestCase {
                principal: 500_000.0,
                ir: 5.0,
                years: 35,
                expected: "2523.44",
                interest_paid_expected: "559844.12",
            },
            TestCase {
                principal: 25_321_323.0,
                ir: 1.2,
                years: 20,
                expected: "118724.61",
                interest_paid_expected: "3172582.59",
            },
        );

        for test_case in test_cases.iter() {
            let monthly_payment = calculate_mortgage(test_case.principal, test_case.ir, test_case.years);
            assert_eq!(format!("{:.2}", monthly_payment), test_case.expected);
            let interest_paid_schedule = calculate_interest_paid_schedule(test_case.principal, monthly_payment, test_case.ir);
            let interest_paid: f64 = interest_paid_schedule.iter().sum();
            assert_eq!(format!("{:.2}", interest_paid), test_case.interest_paid_expected);
        }
    }

    #[test]
    fn test_parse_mortgage() {
        let mortgage_config = parse_mortgage("./sample/mortgage.yaml").unwrap();
        let expected = MortgageConfig {
            mortgages: vec!(
                Mortgage {
                    total: 50_000.0,
                    mortgage_first_payment: Utc.datetime_from_str(
                        "2022-10-10 08:00:00",
                        "%Y-%m-%d %H:%M:%S"
                    ).unwrap(),
                    mortgage_duration_in_years: 25,
                    interest_rate_percentage: 2.6,
                    downpayments: vec!(
                        Downpayment {
                            date: Utc.datetime_from_str(
                                "2021-10-10 08:00:00",
                                "%Y-%m-%d %H:%M:%S"
                            ).unwrap(),
                            sum: 1000.0,
                        },
                        Downpayment {
                            date: Utc.datetime_from_str(
                                "2021-12-12 08:00:00",
                                "%Y-%m-%d %H:%M:%S"
                            ).unwrap(),
                            sum: 20000.0,
                        },
                    ),
                },
            ),
        };
        assert_eq!(mortgage_config, expected);
    }
}
