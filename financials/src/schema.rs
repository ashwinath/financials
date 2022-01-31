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
