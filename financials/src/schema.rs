table! {
    assets (id) {
        id -> Int4,
        transaction_date -> Timestamptz,
        #[sql_name = "type"]
        type_ -> Text,
        amount -> Float8,
    }
}

table! {
    exchange_rates (id) {
        id -> Int4,
        trade_date -> Timestamptz,
        symbol -> Text,
        price -> Float8,
    }
}

table! {
    expenses (id) {
        id -> Int4,
        transaction_date -> Timestamptz,
        #[sql_name = "type"]
        type_ -> Text,
        amount -> Float8,
    }
}

table! {
    incomes (id) {
        id -> Int4,
        transaction_date -> Timestamptz,
        #[sql_name = "type"]
        type_ -> Text,
        amount -> Float8,
    }
}

table! {
    portfolios (id) {
        id -> Text,
        trade_date -> Timestamptz,
        symbol -> Text,
        principal -> Float8,
        nav -> Float8,
        simple_returns -> Float8,
        quantity -> Float8,
    }
}

table! {
    stocks (id) {
        id -> Int4,
        trade_date -> Timestamptz,
        symbol -> Text,
        price -> Float8,
    }
}

table! {
    symbols (id) {
        id -> Int4,
        symbol_type -> Text,
        symbol -> Text,
        base_currency -> Nullable<Text>,
        last_processed_date -> Nullable<Timestamptz>,
    }
}

table! {
    trades (id) {
        id -> Int4,
        date_purchased -> Timestamptz,
        symbol -> Text,
        price_each -> Float8,
        quantity -> Float8,
        trade_type -> Text,
    }
}

allow_tables_to_appear_in_same_query!(
    assets,
    exchange_rates,
    expenses,
    incomes,
    portfolios,
    stocks,
    symbols,
    trades,
);
