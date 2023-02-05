local grafana = import 'grafonnet/grafana.libsonnet';
local graphPanel = grafana.graphPanel;
local dashboard = grafana.dashboard;
local pieChartPanel = grafana.pieChartPanel;
local singlestat = grafana.singlestat;
local tablePanel = grafana.tablePanel;
local row = grafana.row;
local sql = grafana.sql;

local portfolioPieChart = pieChartPanel.new(
  "Portfolio",
  datasource='PostgreSQL'
)
.addTarget(
  sql.target(
    "SELECT
      trade_date as \"time\",
      symbol,
      sum(nav) as \"nav\"
    FROM portfolios
    WHERE
      trade_date=DATE_TRUNC('day', CURRENT_TIMESTAMP - INTERVAL '1 day') + INTERVAL '8 hours'
    group by trade_date, symbol
    order by trade_date",
  )
) + {
  "type": "piechart",
  "options": {
    "legend": {
      "placement": "right",
      "values": [
        "value",
      ],
    }
  },
  "fieldConfig": {
    "defaults": {
      "unit": "currencyUSD",
    },
  },
}; // hack to use the new piechart

local currentStateTable = tablePanel.new(
  "Current Portfolio",
  datasource='PostgreSQL',
  styles=[
    {
      "type": "hidden",
      "pattern": "time",
    },
    {
      "pattern": "symbol",
      "alias": "Symbol",
    },
    {
      "unit": "currencyUSD",
      "type": "number",
      "alias": "Current Price (USD)",
      "decimals": 2,
      "pattern": "current_price",
    },
    {
      "pattern": "quantity",
      "alias": "Quantity",
      "decimals": 0,
      "type": "number",
    },
    {
      "unit": "currencyUSD",
      "type": "number",
      "alias": "NAV",
      "decimals": 2,
      "pattern": "nav",
    },
    {
      "unit": "currencyUSD",
      "type": "number",
      "alias": "Principal",
      "decimals": 2,
      "pattern": "principal",
    },
    {
      "unit": "percentunit",
      "type": "number",
      "alias": "Returns",
      "decimals": 2,
      "pattern": "returns",
    },
    {
      "unit": "percentunit",
      "type": "number",
      "alias": "Weight",
      "decimals": 2,
      "pattern": "percentage",
    },
  ],
)
.addTarget(
  sql.target("
    WITH total AS (
      SELECT
        sum(nav) AS total
      FROM portfolios
      WHERE
        trade_date=DATE_TRUNC('day', CURRENT_TIMESTAMP - INTERVAL '1 day') + INTERVAL '8 hours'
    ),
    stock as (
      SELECT
        symbol,
        price
      FROM stocks
      WHERE
        trade_date = (
          SELECT
            trade_date
          FROM stocks
          ORDER BY trade_date desc
          LIMIT 1
        )
    )
    SELECT
      trade_date as time,
      stock.symbol as symbol,
      stock.price as current_price,
      quantity as quantity,
      sum(nav) as nav,
      sum(principal) as principal,
      (sum(nav) - sum(principal)) / sum(principal) as returns,
      sum(nav)/total.total as percentage
    FROM portfolios inner join stock on portfolios.symbol = stock.symbol, total
    WHERE
      trade_date=DATE_TRUNC('day', CURRENT_TIMESTAMP - INTERVAL '1 day') + INTERVAL '8 hours'
    group by trade_date, stock.symbol, total.total, quantity, stock.price
    order by trade_date",
    format='table',
  )
);

local createStat(name, unit, query) =
  singlestat.new(
    name,
    format=unit,
    datasource='PostgreSQL',
  )
  .addTarget(
    sql.target(
      query,
      format='time_series',
    )
  );

local currentPrincipal = createStat(
  name="Principal",
  unit="currencyUSD",
  query="SELECT
    trade_date as \"time\",
    sum(principal) as \"principal\"
  FROM portfolios
  WHERE
    trade_date=DATE_TRUNC('day', CURRENT_TIMESTAMP - INTERVAL '1 day') + INTERVAL '8 hours'
  group by trade_date
  order by trade_date",
);

local currentNAV = createStat(
  name="NAV",
  unit="currencyUSD",
  query="SELECT
    trade_date as \"time\",
    sum(nav) as \"nav\"
  FROM portfolios
  WHERE
    trade_date=DATE_TRUNC('day', CURRENT_TIMESTAMP - INTERVAL '1 day') + INTERVAL '8 hours'
  group by trade_date
  order by trade_date",
);

local currentSimpleReturns = createStat(
  name="Simple Returns",
  unit="percentunit",
  query="SELECT
    trade_date as \"time\",
    (sum(nav) - sum(principal)) / sum(principal) as \"returns\"
  FROM portfolios
  WHERE
    trade_date=DATE_TRUNC('day', CURRENT_TIMESTAMP - INTERVAL '1 day') + INTERVAL '8 hours'
  group by trade_date
  order by trade_date",
);

local createPanel(name, unit, query, legend_show, stack=false, points=false) =
  graphPanel.new(
    name,
    datasource='PostgreSQL',
    format=unit,
    legend_show=legend_show,
    fillGradient=4,
    linewidth=2,
    stack=stack,
    nullPointMode='null as zero',
    points=points,
    pointradius=5
  )
  .addTarget(
    sql.target(
      query,
      format='time_series',
    )
  );

local portfolioSimpleReturns = createPanel(
  name='Portfolio Simple Returns',
  unit='percentunit',
  query='SELECT
    trade_date as "time",
    (sum(nav) - sum(principal)) / sum(principal) as "returns"
  FROM portfolios
  WHERE
    $__timeFilter(trade_date)
  group by trade_date
  order by trade_date',
  legend_show=false,
);

local portfolioNAV = createPanel(
  name='Portfolio NAV',
  unit='currencyUSD', // it's actually SGD but grafana wants USD
  query='SELECT
    trade_date as "time",
    sum(nav) as "nav",
    sum(principal) as "principal"
  FROM portfolios
  WHERE
    $__timeFilter(trade_date)
  group by trade_date
  order by trade_date',
  legend_show=true,
);

local nav = createPanel(
  name='NAV',
  unit='currencyUSD', // it's actually SGD but grafana wants USD
  query='SELECT
    trade_date as "time",
    symbol,
    sum(nav) as "nav"
  FROM portfolios
  WHERE
    $__timeFilter(trade_date)
  group by trade_date, symbol
  order by trade_date',
  legend_show=true,
);

local simpleReturns = createPanel(
  name='Simple Returns',
  unit='percentunit',
  query='SELECT
    trade_date as "time",
    symbol,
    (sum(nav) - sum(principal))/ sum(principal) as "returns"
  FROM portfolios
  WHERE
    $__timeFilter(trade_date)
  group by trade_date, symbol
  order by trade_date',
  legend_show=true,
);

// Summary of Financials
local netAssets = createPanel(
  name='Net Assets',
  unit='currencyUSD',
  query='SELECT
    transaction_date as "time",
    SUM(amount) as "Amount"
  FROM assets
  WHERE
    $__timeFilter(transaction_date)
  group by transaction_date
  order by transaction_date',
  legend_show=false,
  stack=true,
);

local liquidAssets = createPanel(
  name='Net Liquid Assets',
  unit='currencyUSD',
  query='SELECT
    transaction_date as "time",
    type,
    amount AS "Amount"
  FROM assets
  WHERE
    $__timeFilter(transaction_date)
    AND type IN (\'Bank\', \'Investments\', \'Bonds\')
  group by transaction_date, type, amount
  order by transaction_date',
  legend_show=true,
  stack=true,
);

local nonLiquidAssets = createPanel(
  name='Net Illiquid Assets',
  unit='currencyUSD',
  query='SELECT
    transaction_date as "time",
    type,
    amount AS "Amount"
  FROM assets
  WHERE
    $__timeFilter(transaction_date)
    AND type IN (\'OA\', \'SA\', \'Medisave\', \'SRS\', \'House\')
  group by transaction_date, type, amount
  order by transaction_date',
  legend_show=true,
  stack=true,
);

local salary = createPanel(
  name='Salary',
  unit='currencyUSD',
  query='SELECT
    transaction_date as "time",
    type,
    amount AS "Amount"
  FROM incomes
  WHERE
    $__timeFilter(transaction_date)
  group by transaction_date, type, amount
  order by transaction_date',
  legend_show=true,
  stack=true,
);

local expenses = createPanel(
  name='Expenses',
  unit='currencyUSD',
  query='WITH cc as (
    SELECT
      transaction_date as "transaction_date",
      sum(amount) as "amount"
    FROM
      expenses
    WHERE
      $__timeFilter(transaction_date)
      AND type in (\'Credit Card\', \'Reimbursement\')
    group by transaction_date
    order by transaction_date
  )
  SELECT
      transaction_date as "time",
      type,
      amount AS "Amount"
  FROM expenses
  WHERE
      $__timeFilter(transaction_date)
      AND type not in (\'Credit Card\', \'Reimbursement\')
      AND type not like \'Special:%\'
  UNION select transaction_date as "time", \'Credit Card\' as type, amount from cc
  group by transaction_date, type, amount
  order by time',
  legend_show=true,
  stack=true,
);

local special_expenses = createPanel(
  name='Special Expenses',
  unit='currencyUSD',
  query='SELECT
      transaction_date as "time",
      type,
      amount AS "Amount"
  FROM expenses
  WHERE
      $__timeFilter(transaction_date)
      AND type like \'Special:%\'
  group by transaction_date, type, amount
  order by time',
  legend_show=true,
  stack=true,
  points=true,
);

local savingsRate = createPanel(
  name='Savings Rate',
  unit='percentunit',
  query='WITH expenditure AS (
  select date_trunc(\'month\', transaction_date) as time, sum(amount) as Amount from expenses where $__timeFilter(transaction_date) and type not like \'Special:%\' group by time
),
income as (
  select date_trunc(\'month\', transaction_date) as time, sum(amount) as Amount from incomes where $__timeFilter(transaction_date) and type in (\'Base\', \'Base Bonus\') group by time
)
SELECT
  expenditure.time + interval \'1 month\' - interval \'1 day\' AS "time",
  (1 - (expenditure.amount / income.amount)) as \"Savings Rate\"
FROM expenditure inner join income on income.time = expenditure.time
WHERE $__timeFilter(expenditure.time)
ORDER BY expenditure.time;',
  legend_show=true,
  stack=true,
);

local emergencyFunds = createPanel(
  name='Emergency Funds',
  unit='Months',
  query="WITH avg_exp AS (
    select expense_date + INTERVAL '1 day' as \"expense_date\", amount from average_expenditures
),
bank AS (
    select transaction_date, sum(amount) as amount from assets where type in ('Bank', 'Bonds') group by transaction_date
)
select
    a.expense_date as \"time\",
    b.amount / a.amount as \"months\"
from
    avg_exp a inner join bank b on a.expense_date = b.transaction_date
WHERE $__timeFilter(a.expense_date)
ORDER BY a.expense_date;",
  legend_show=false,
  stack=true,
);

local runway = createPanel(
  name='Runway Based on 70% Equity + Bank + Bonds',
  unit='Months',
  query="WITH avg_exp AS (
    select expense_date + INTERVAL '1 day' as \"expense_date\", amount from average_expenditures
),
bank AS (
    select transaction_date, sum(amount) as amount from assets where type in ('Bank', 'Bonds') group by transaction_date
),
equity AS (
    select transaction_date, amount from assets where type = 'Investments'
)
select
    a.expense_date as \"time\",
    (b.amount + (c.amount * 0.7)) / a.amount as \"months\"
from
    avg_exp a inner join bank b on a.expense_date = b.transaction_date
    inner join equity c on a.expense_date = c.transaction_date
WHERE $__timeFilter(a.expense_date)
ORDER BY a.expense_date;",
  legend_show=false,
  stack=true,
);

local fiQuotient = createPanel(
  name='Financial Independence Quotient (3% Withdrawal)',
  unit='percentunit',
  query="WITH avg_exp AS (
    select expense_date + INTERVAL '1 day' as \"expense_date\", amount from average_expenditures
),
equity AS (
    select transaction_date, sum(amount) as amount from assets where type in ('Investments', 'Bonds') group by transaction_date
)
select
    a.expense_date as \"time\",
    (0.03 * c.amount) / (a.amount * 12) as \"quotient\"
from
    avg_exp a inner join equity c on a.expense_date = c.transaction_date
WHERE $__timeFilter(a.expense_date)
ORDER BY a.expense_date;",
  legend_show=false,
  stack=true,
);

local mortgage_payment = createPanel(
  name='Mortgage Payment',
  unit='currencyUSD',
  query='SELECT
    date as "time",
    interest_paid as "Interest Paid",
    principal_paid AS "Principal Paid"
FROM mortgage
WHERE
    $__timeFilter(date)
group by date, interest_paid, principal_paid
order by time',
  legend_show=true,
  stack=true,
);

local mortgage_serviced = createPanel(
  name='Mortgage Serviced',
  unit='currencyUSD',
  query='SELECT
    date as "time",
    total_interest_paid as "Total Interest Paid",
    total_principal_paid AS "Total Principal Paid"
FROM mortgage
WHERE
    $__timeFilter(date)
group by date, total_interest_paid, total_principal_paid
order by time',
  legend_show=true,
  stack=true,
);

local mortgage_summary = createPanel(
  name='Mortgage Summary',
  unit='currencyUSD',
  query='SELECT
    date as "time",
    total_interest_left as "Interest Left",
    total_principal_left AS "Principal Left"
FROM mortgage
WHERE
    $__timeFilter(date)
group by date, total_interest_left, total_principal_left
order by time',
  legend_show=true,
  stack=true,
);

dashboard.new(
  'Financials',
  schemaVersion=16,
  tags=['financials'],
  time_from='now-1y',
  editable=true,
  graphTooltip='shared_tooltip',
)
// Summary Of Financials
.addPanel(
  row.new(
    title="Financials Summary"
  ),
  gridPos={ h: 1, w: 12, x: 0, y: 0 },
)
# Assets
.addPanel(
  netAssets,
  gridPos={ h: 8, w: 8, x: 0, y: 1 },
)
.addPanel(
  liquidAssets,
  gridPos={ h: 8, w: 8, x: 8, y: 1 },
)
.addPanel(
  nonLiquidAssets,
  gridPos={ h: 8, w: 8, x: 16, y: 1 },
)

# Salary/Expenses Information
.addPanel(
  row.new(
    title="Inflow/Outflow"
  ),
  gridPos={ h: 1, w: 12, x: 0, y: 9 },
)
.addPanel(
  salary,
  gridPos={ h: 8, w: 8, x: 0, y: 10 },
)
.addPanel(
  expenses,
  gridPos={ h: 8, w: 8, x: 8, y: 10 },
)
.addPanel(
  savingsRate,
  gridPos={ h: 8, w: 8, x: 16, y: 10 },
)
.addPanel(
  special_expenses,
  gridPos={ h: 8, w: 8, x: 0, y: 18 },
)

// Mortgage
.addPanel(
  row.new(
    title="Mortgage"
  ),
  gridPos={ h: 1, w: 12, x: 0, y: 19 },
)
.addPanel(
  mortgage_summary,
  gridPos={ h: 8, w: 8, x: 8, y: 20 },
)
.addPanel(
  mortgage_serviced,
  gridPos={ h: 8, w: 8, x: 16, y: 20 },
)
.addPanel(
  mortgage_payment,
  gridPos={ h: 8, w: 8, x: 0, y: 28 },
)

# FI Ratios
.addPanel(
  row.new(
    title="Financial Health"
  ),
  gridPos={ h: 1, w: 12, x: 0, y: 29 },
)
.addPanel(
  emergencyFunds,
  gridPos={ h: 8, w: 8, x: 0, y: 30 },
)
.addPanel(
  runway,
  gridPos={ h: 8, w: 8, x: 8, y: 30 },
)
.addPanel(
  fiQuotient,
  gridPos={ h: 8, w: 8, x: 16, y: 30 },
)

// Current State
.addPanel(
  row.new(
    title="Current Investments"
  ),
  gridPos={ h: 1, w: 12, x: 0, y: 31 },
)
.addPanel(
  currentStateTable,
  gridPos={ h: 8, w: 10, x: 0, y: 32 },
)
.addPanel(
  portfolioPieChart,
  gridPos={ h: 8, w: 6, x: 10, y: 32 },
)
.addPanel(
  currentNAV,
  gridPos={ h: 4, w: 4, x: 16, y: 32 },
)
.addPanel(
  currentPrincipal,
  gridPos={ h: 4, w: 4, x: 20, y: 32 },
)
.addPanel(
  currentSimpleReturns,
  gridPos={ h: 4, w: 4, x: 16, y: 40 },
)

// PERFORMANCE
.addPanel(
  row.new(
    title="Investment Historical Performance"
  ),
  gridPos={ h: 1, w: 12, x: 0, y: 41 },
)
.addPanel(
  portfolioSimpleReturns,
  gridPos={ h: 8, w: 12, x: 0, y: 42 },
)
.addPanel(
  portfolioNAV,
  gridPos={ h: 8, w: 12, x: 12, y: 42 },
)
.addPanel(
  simpleReturns,
  gridPos={ h: 8, w: 12, x: 0, y: 50 },
)
.addPanel(
  nav,
  gridPos={ h: 8, w: 12, x: 12, y: 50 },
)
