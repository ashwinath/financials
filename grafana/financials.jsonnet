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

local createPanel(name, unit, query, legend_show, stack=false) =
  graphPanel.new(
    name,
    datasource='PostgreSQL',
    format=unit,
    legend_show=legend_show,
    fillGradient=4,
    linewidth=2,
    stack=stack,
    nullPointMode='null as zero'
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
    type,
    amount
  FROM assets
  WHERE
    $__timeFilter(transaction_date)
  group by transaction_date,type, amount
  order by transaction_date',
  legend_show=false,
  stack=true,
);

local salary = createPanel(
  name='Salary',
  unit='currencyUSD',
  query='SELECT
    transaction_date as "time",
    type,
    amount
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
      amount
  FROM expenses
  WHERE
      $__timeFilter(transaction_date)
      AND type not in (\'Credit Card\', \'Reimbursement\')
  UNION select transaction_date as "time", \'Credit Card\' as type, amount from cc
  group by transaction_date, type, amount
  order by time',
  legend_show=true,
  stack=true,
);

local savingsRate = createPanel(
  name='Savings Rate',
  unit='percentunit',
  query='WITH expenditure AS (
  select date_trunc(\'month\', transaction_date) as time, sum(amount) as amount from expenses where $__timeFilter(transaction_date) group by time
),
income as (
  select date_trunc(\'month\', transaction_date) as time, sum(amount) as amount from incomes where $__timeFilter(transaction_date) and type in (\'base\', \'base_bonus\') group by time
)
SELECT
  expenditure.time + interval \'1 month\' - interval \'1 day\' AS "time",
  (1 - (expenditure.amount / income.amount)) as savings_rate
FROM expenditure inner join income on income.time = expenditure.time
WHERE $__timeFilter(expenditure.time)
ORDER BY expenditure.time;',
  legend_show=true,
  stack=true,
);

dashboard.new(
  'Financials',
  schemaVersion=16,
  tags=['financials'],
  time_from='now-90d',
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
.addPanel(
  netAssets,
  gridPos={ h: 8, w: 12, x: 0, y: 1 },
)
.addPanel(
  salary,
  gridPos={ h: 8, w: 12, x: 12, y: 1 },
)
.addPanel(
  expenses,
  gridPos={ h: 8, w: 12, x: 0, y: 9 },
)
.addPanel(
  savingsRate,
  gridPos={ h: 8, w: 12, x: 12, y: 1 },
)
// CURRENT STATE
.addPanel(
  row.new(
    title="Current Investments"
  ),
  gridPos={ h: 1, w: 12, x: 0, y: 17 },
)
.addPanel(
  currentStateTable,
  gridPos={ h: 8, w: 10, x: 0, y: 18 },
)
.addPanel(
  portfolioPieChart,
  gridPos={ h: 8, w: 6, x: 10, y: 18 },
)
.addPanel(
  currentNAV,
  gridPos={ h: 4, w: 4, x: 16, y: 18 },
)
.addPanel(
  currentPrincipal,
  gridPos={ h: 4, w: 4, x: 20, y: 18 },
)
.addPanel(
  currentSimpleReturns,
  gridPos={ h: 4, w: 4, x: 16, y: 22 },
)

// PERFORMANCE
.addPanel(
  row.new(
    title="Investment Historical Performance"
  ),
  gridPos={ h: 1, w: 12, x: 0, y: 23 },
)
.addPanel(
  portfolioSimpleReturns,
  gridPos={ h: 8, w: 12, x: 0, y: 24 },
)
.addPanel(
  portfolioNAV,
  gridPos={ h: 8, w: 12, x: 12, y: 24 },
)
.addPanel(
  simpleReturns,
  gridPos={ h: 8, w: 12, x: 0, y: 32 },
)
.addPanel(
  nav,
  gridPos={ h: 8, w: 12, x: 12, y: 32 },
)
