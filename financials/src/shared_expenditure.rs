use chrono::{DateTime, Utc};
use std::error::Error;
use diesel::pg::PgConnection;
use diesel::dsl::sum;
use diesel::{insert_into, ExpressionMethods, QueryDsl, RunQueryDsl, TextExpressionMethods};
use crate::schema::expenses::dsl::expenses;
use crate::schema::shared_expense::dsl::shared_expense;
use crate::models::SharedExpenseWithId;
// Needs to go inside expenses as a row called shared_expenses
// Needs to go inside special expenses as a row called special something 
// Average expense needs to be recalculated.
// Order of execution is important

struct DateAmountPair {
    pub amount: f64,
}

pub fn populate_shared_expenditure(conn: &mut PgConnection) -> Result<(), Box<dyn Error>> {
    // Handle non special expense
    //let non_special_expenses = shared_expense
        //.select((
            //crate::schema::shared_expense::expense_date,
            //sum(crate::schema::shared_expense::amount)
        //))
        //.filter(crate::schema::shared_expense::type_.not_like(String::from("Special:%")))
        //.group_by(crate::schema::shared_expense::expense_date)
        //.load::<(DateTime<Utc>, f64)>(conn)?;

    // Handle special expense

    Ok(())
}
