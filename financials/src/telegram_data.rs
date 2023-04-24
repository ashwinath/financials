use chrono::{DateTime, Datelike, TimeZone, Utc};
use chronoutil::delta::shift_months;
use diesel::pg::PgConnection;
use reqwest;
use serde::Deserialize;
use std::fmt;

use crate::models::{Expense, SharedExpense};
use crate::schema::expenses::dsl::expenses;
use crate::schema::shared_expense::dsl::shared_expense;
use diesel::{insert_into, RunQueryDsl};
use std::error::Error;

// This is hard coded.
// I do not see this value changing as this as this feature was created only then.
const START_MONTH: u32 = 3;
const START_YEAR: i32 = 2023;

#[derive(Debug, Clone)]
pub struct TelegramBotError {
    pub message: String,
}

impl Error for TelegramBotError {}

impl fmt::Display for TelegramBotError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Telegram Bot API Error")
    }
}

pub fn load_telegram_data(
    conn: &mut PgConnection,
    telegram_bot_endpoint: &str,
) -> Result<(), Box<dyn Error>> {
    let expense_responses = download_expenses(telegram_bot_endpoint, "expenses")?;
    let expense_list: Vec<Expense> = expense_responses
        .iter()
        .map(map_expense_response_to_expense)
        .collect();
    insert_into(expenses).values(&expense_list).execute(conn)?;

    let shared_expense_responses = download_expenses(telegram_bot_endpoint, "shared-expenses")?;
    let shared_expense_list: Vec<SharedExpense> = shared_expense_responses
        .iter()
        .map(map_expense_response_to_shared_expense)
        .collect();
    insert_into(shared_expense)
        .values(&shared_expense_list)
        .execute(conn)?;

    Ok(())
}

fn map_expense_response_to_shared_expense(e: &ExpenseStruct) -> SharedExpense {
    let type_ = e.type_.clone();
    return SharedExpense {
        expense_date: e.date,
        type_,
        amount: e.amount,
    };
}

fn map_expense_response_to_expense(e: &ExpenseStruct) -> Expense {
    let type_ = e.type_.clone();
    return Expense {
        id: None,
        transaction_date: e.date,
        type_,
        amount: e.amount,
    };
}

#[derive(Debug, Deserialize)]
struct ExpenseResponse {
    pub expenses: Vec<ExpenseStruct>,
}

#[derive(Debug, Deserialize)]
struct ExpenseStruct {
    pub date: DateTime<Utc>,
    #[serde(rename(deserialize = "type"))]
    pub type_: String,
    pub amount: f64,
}

fn download_expenses(
    telegram_bot_endpoint: &str,
    expense_type: &str,
) -> Result<Vec<ExpenseStruct>, Box<dyn Error>> {
    let mut processing_date = Utc
        .with_ymd_and_hms(START_YEAR, START_MONTH, 1, 8, 0, 0)
        .unwrap();
    let today = Utc::now();

    let mut res: Vec<ExpenseStruct> = Vec::new();
    while processing_date < today {
        let url = format!(
            "{}/{}?month={:02}&year={:04}",
            telegram_bot_endpoint,
            expense_type,
            processing_date.month(),
            processing_date.year(),
        );

        processing_date = shift_months(processing_date, 1);
        let mut result: ExpenseResponse = call_telegram_bot(&url)?;
        res.append(&mut result.expenses);
    }

    Ok(res)
}

fn call_telegram_bot<T>(url: &str) -> Result<T, Box<dyn Error>>
where
    T: for<'de> serde::Deserialize<'de>,
{
    let response = reqwest::blocking::get(url)?;

    match response.status() {
        reqwest::StatusCode::OK => return Ok(response.json::<T>()?),
        e => {
            return Err(TelegramBotError {
                message: format!(
                    "status code: {} {}",
                    e.as_str(),
                    e.canonical_reason().unwrap()
                ),
            }
            .into())
        }
    }
}
