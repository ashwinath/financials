table! {
    assets (id) {
        id -> Int4,
        transaction_date -> Timestamptz,
        #[sql_name = "type"]
        type_ -> Text,
        amount -> Float8,
        created_at -> Timestamptz,
        updated_at -> Timestamptz,
    }
}

table! {
    expenses (id) {
        id -> Int4,
        transaction_date -> Timestamptz,
        #[sql_name = "type"]
        type_ -> Text,
        amount -> Float8,
        created_at -> Timestamptz,
        updated_at -> Timestamptz,
    }
}

table! {
    incomes (id) {
        id -> Int4,
        transaction_date -> Timestamptz,
        #[sql_name = "type"]
        type_ -> Text,
        amount -> Float8,
        created_at -> Timestamptz,
        updated_at -> Timestamptz,
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
        created_at -> Timestamptz,
        updated_at -> Timestamptz,
    }
}

allow_tables_to_appear_in_same_query!(
    assets,
    expenses,
    incomes,
    trades,
);
