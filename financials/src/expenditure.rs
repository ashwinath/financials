use chrono::{DateTime, Utc, Duration};
use chronoutil::delta::shift_months;
use crate::schema::average_expenditures::dsl::average_expenditures;
use std::error::Error;
use crate::schema::expenses::dsl::expenses;
use crate::models::AverageExpenditure;
use diesel::dsl::sum;
use diesel::{insert_into, ExpressionMethods, QueryDsl, RunQueryDsl, TextExpressionMethods};
use diesel::pg::PgConnection;
use diesel::pg::upsert::excluded;

use crate::utils::{get_last_day_of_month, is_last_day_of_month};

const WINDOW_PERIOD: i32 = 6;

// Calculates monthly expenditure rates based on half yearly rolling window.
// Expenditure should not include taxes as if we retire there is no tax.
pub fn calculate_average_expenditure(conn: &mut PgConnection) -> Result<(), Box<dyn Error>> {
    let first_date = expenses
        .select(crate::schema::expenses::transaction_date)
        .order(crate::schema::expenses::transaction_date.asc())
        .first::<DateTime<Utc>>(conn)?;

    let mut current_date = if is_last_day_of_month(first_date) {
        first_date.clone()
    } else {
        shift_months(first_date, 1)
    };

    // Get last day of month
    current_date = get_last_day_of_month(current_date);

    // Start from 1 year later.
    current_date = shift_months(current_date, WINDOW_PERIOD);
    let tomorrow = chrono::offset::Utc::now() + Duration::days(1);

    let mut all_average_expenditures: Vec<AverageExpenditure> = Vec::new();
    while current_date < tomorrow {
        let yearly_expenditure = expenses
            .select(sum(crate::schema::expenses::amount))
            .filter(crate::schema::expenses::transaction_date.gt(shift_months(current_date, -WINDOW_PERIOD)))
            .filter(crate::schema::expenses::transaction_date.le(current_date))
            .filter(crate::schema::expenses::type_.ne(String::from("Tax")))
            .filter(crate::schema::expenses::type_.not_like(String::from("Special:%")))
            .first::<Option<f64>>(conn)?;

        if let Some(i) = yearly_expenditure {
            let avg_expenditure = AverageExpenditure {
                id: None,
                expense_date: current_date.clone(),
                amount: i / f64::from(WINDOW_PERIOD),
            };
            all_average_expenditures.push(avg_expenditure);
        }

        current_date = get_last_day_of_month(shift_months(current_date, 1));
    }
    insert_into(average_expenditures)
        .values(all_average_expenditures)
        .on_conflict(crate::schema::average_expenditures::expense_date)
        .do_update()
        .set(
            crate::schema::average_expenditures::amount.eq(
                excluded(crate::schema::average_expenditures::amount)
            )
        )
        .execute(conn)?;

    Ok(())
}
