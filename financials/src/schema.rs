// @generated automatically by Diesel CLI.

diesel::table! {
    assets (id) {
        id -> Int4,
        transaction_date -> Timestamptz,
        #[sql_name = "type"]
        type_ -> Text,
        amount -> Float8,
    }
}

diesel::table! {
    average_expenditures (id) {
        id -> Int4,
        expense_date -> Timestamptz,
        amount -> Float8,
    }
}

diesel::table! {
    exchange_rates (id) {
        id -> Int4,
        trade_date -> Timestamptz,
        symbol -> Text,
        price -> Float8,
    }
}

diesel::table! {
    expenses (id) {
        id -> Int4,
        transaction_date -> Timestamptz,
        #[sql_name = "type"]
        type_ -> Text,
        amount -> Float8,
    }
}

diesel::table! {
    incomes (id) {
        id -> Int4,
        transaction_date -> Timestamptz,
        #[sql_name = "type"]
        type_ -> Text,
        amount -> Float8,
    }
}

diesel::table! {
    mortgage (id) {
        id -> Int4,
        date -> Timestamptz,
        interest_paid -> Float8,
        principal_paid -> Float8,
        total_principal_paid -> Float8,
        total_interest_paid -> Float8,
        total_principal_left -> Float8,
        total_interest_left -> Float8,
    }
}

diesel::table! {
    portfolios (id) {
        id -> Int4,
        trade_date -> Timestamptz,
        symbol -> Text,
        principal -> Float8,
        nav -> Float8,
        simple_returns -> Float8,
        quantity -> Float8,
    }
}

diesel::table! {
    shared_expense (id) {
        id -> Int4,
        expense_date -> Timestamptz,
        #[sql_name = "type"]
        type_ -> Text,
        amount -> Float8,
    }
}

diesel::table! {
    stocks (id) {
        id -> Int4,
        trade_date -> Timestamptz,
        symbol -> Text,
        price -> Float8,
    }
}

diesel::table! {
    symbols (id) {
        id -> Int4,
        symbol_type -> Text,
        symbol -> Text,
        base_currency -> Nullable<Text>,
        last_processed_date -> Nullable<Timestamptz>,
    }
}

diesel::table! {
    trades (id) {
        id -> Int4,
        date_purchased -> Timestamptz,
        symbol -> Text,
        price_each -> Float8,
        quantity -> Float8,
        trade_type -> Text,
    }
}

diesel::allow_tables_to_appear_in_same_query!(
    assets,
    average_expenditures,
    exchange_rates,
    expenses,
    incomes,
    mortgage,
    portfolios,
    shared_expense,
    stocks,
    symbols,
    trades,
);
