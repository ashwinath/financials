use chrono::{DateTime, Utc};
use std::error::Error;
use diesel::pg::PgConnection;
use diesel::pg::sql_types::Timestamptz;
use diesel::sql_types::Double;
use diesel::{insert_into, RunQueryDsl};
use diesel::dsl::sql_query;

use crate::models::Expense;
use crate::schema::expenses::dsl::expenses;
// Needs to go inside expenses as a row called shared_expenses
// Needs to go inside special expenses as a row called special something 
// Average expense needs to be recalculated.
// Order of execution is important

const NON_SPECIAL_SHARED_EXPENSE_TYPE: &str = "Shared Expense";
const NON_SPECIAL_EXPENSES_SQL: &str = "SELECT expense_date, sum(amount) AS total FROM shared_expense WHERE type NOT LIKE 'Special:%' GROUP BY expense_date";

const SPECIAL_SHARED_EXPENSE_TYPE: &str = "Special:Shared Expense";
const SPECIAL_EXPENSES_SQL: &str = "SELECT expense_date, sum(amount) AS total FROM shared_expense WHERE type LIKE 'Special:%' GROUP BY expense_date";

#[derive(Debug, QueryableByName)]
struct DateAmountPair {
    #[diesel(sql_type = Timestamptz)]
    pub expense_date: DateTime<Utc>,
    #[diesel(sql_type = Double)]
    pub total: f64,
}

pub fn populate_shared_expenditure(conn: &mut PgConnection) -> Result<(), Box<dyn Error>> {
    // Handle non special expense
    insert_expense(conn, NON_SPECIAL_EXPENSES_SQL, NON_SPECIAL_SHARED_EXPENSE_TYPE)?;

    // Handle special expense
    insert_expense(conn, SPECIAL_EXPENSES_SQL, SPECIAL_SHARED_EXPENSE_TYPE)?;

    Ok(())
}

fn insert_expense(conn: &mut PgConnection, sql: &str, expense_type: &str) -> Result<(), Box<dyn Error>> {
    let expenses_queries: Vec<DateAmountPair> = sql_query(sql).load(conn)?;
    let non_special_expenses: Vec<Expense> = expenses_queries.iter().map(|x| Expense {
        id: None,
        transaction_date: x.expense_date,
        type_: String::from(expense_type),
        amount: x.total,
    }).collect();

    insert_into(expenses)
        .values(non_special_expenses)
        .execute(conn)?;

    Ok(())
}
