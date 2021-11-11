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

local createPanel(name, unit, query, legend_show) =
  graphPanel.new(
    name,
    datasource='PostgreSQL',
    format=unit,
    legend_show=legend_show,
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
    (sum(nav) - sum(principal)) / sum(principal)
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

dashboard.new(
  'Financials',
  schemaVersion=16,
  tags=['financials'],
  time_from='now-90d',
  editable=true,
)
// CURRENT STATE
.addPanel(
  row.new(
    title="Current State"
  ),
  gridPos= { h: 1, w: 12, x: 0, y: 0 },
)
.addPanel(
  currentStateTable,
  gridPos= { h: 8, w: 10, x: 0, y: 0 },
)
.addPanel(
  portfolioPieChart,
  gridPos= { h: 8, w: 6, x: 10, y: 0 },
)
.addPanel(
  currentNAV,
  gridPos= { h: 4, w: 4, x: 16, y: 0 },
)
.addPanel(
  currentPrincipal,
  gridPos= { h: 4, w: 4, x: 20, y: 0 },
)
.addPanel(
  currentSimpleReturns,
  gridPos= { h: 4, w: 4, x: 16, y: 0 },
)
// PERFORMANCE
.addPanel(
  row.new(
    title="Historical Performance"
  ),
  gridPos= { h: 11, w: 12, x: 0, y: 1 },
)
.addPanel(
  portfolioSimpleReturns,
  gridPos= { h: 10, w: 12, x: 0, y: 2 },
)
.addPanel(
  portfolioNAV,
  gridPos= { h: 10, w: 12, x: 12, y: 2 },
)
.addPanel(
  simpleReturns,
  gridPos= { h: 10, w: 12, x: 0, y: 3 },
)
.addPanel(
  nav,
  gridPos= { h: 10, w: 12, x: 12, y: 3 },
)
