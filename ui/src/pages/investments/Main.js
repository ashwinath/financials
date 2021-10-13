import React from 'react';
import moment from 'moment';

import {
  EuiPageTemplate,
  EuiFlexGroup,
  EuiFlexItem,
  EuiHorizontalRule,
  EuiBasicTable,
  EuiSuperSelect,
  EuiSpacer,
} from '@elastic/eui';

import {
  timeFormatter,
  niceTimeFormatByDay,
  Chart,
  Settings,
  AreaSeries,
  Axis,
  Position,
  Partition,
} from '@elastic/charts';

import { EUI_CHARTS_THEME_LIGHT } from '@elastic/eui/dist/eui_charts_theme';

import { SideBar, Stat } from "../../components";
import { LoadingPage } from "../";
import { useLoginHook } from "../../hooks";
import { useDispatch, useSelector } from 'react-redux';

import {
  queryPortfolio,
  updateQueryPeriodInMonths,
} from '../../redux/investmentsSlice';

import { 
  getDateFromPeriod,
  formatMoneyGraph,
  capitaliseAll,
  formatMoney,
  formatPercent,
} from "../../utils";

const columns = [
  {
    field: 'symbol',
    name: 'Symbol',
    truncateText: true,
    render: (field) => capitaliseAll(field),
  },
  {
    field: 'principal',
    name: 'Principal',
    truncateText: true,
    render: (field) => formatMoney(field),
  },
  {
    field: 'nav',
    name: 'NAV',
    truncateText: true,
    render: (field) => formatMoney(field),
  },
  {
    field: 'quantity',
    name: 'Quantity',
    truncateText: true,
  },
  {
    field: 'simple_returns',
    name: 'Returns',
    truncateText: true,
    render: (field) => formatPercent(field),
  },
  {
    field: 'percent',
    name: 'Percentage',
    truncateText: true,
    render: (field) => formatPercent(field),
  },
];

const timeOptions = [
  {
    value: 1,
    inputDisplay: "Past Month",
  },
  {
    value: 3,
    inputDisplay: "Past Quarter",
  },
  {
    value: 6,
    inputDisplay: "Past Half",
  },
  {
    value: 12,
    inputDisplay: "Past Year",
  },
  {
    value: 24,
    inputDisplay: "Past 2 Years",
  },
  {
    value: 36,
    inputDisplay: "Past 3 Years",
  },
  {
    value: 60,
    inputDisplay: "Past 5 Years",
  },
  {
    value: 120,
    inputDisplay: "Past 10 Years",
  },
  {
    value: 240,
    inputDisplay: "Past 20 Years",
  },
  {
    value: 360,
    inputDisplay: "Past 30 Years",
  },
]

export function InvestmentsMainPage() {
  const status = useLoginHook();
  const dispatch = useDispatch();
  const {
    queryPeriodInMonths,
    portfolio,
    portfolioLoaded,
    portfolioLoading,
    shouldReload,
  } = useSelector((state) => state.investments);

  if ((!portfolioLoading && !portfolioLoaded) || shouldReload) {
    dispatch(queryPortfolio(getDateFromPeriod(queryPeriodInMonths)))
  }

  if (status === "loading" || portfolioLoading || portfolio.length === 0) {
    return <LoadingPage/>;
  }

  const allSymbols = {};
  portfolio.forEach((item) => {
    if (item.symbol in allSymbols) {
      allSymbols[item.symbol].push(item);
    } else {
      allSymbols[item.symbol] = [item];
    }
  });

  const aggregatedNavPrincipal = {}; // trade_date: {nav, principal}
  Object.entries(allSymbols).forEach(([symbol, items]) => {
    for (const item of items) {
      if (item.trade_date in aggregatedNavPrincipal) {
        aggregatedNavPrincipal[item.trade_date] = {
          nav: aggregatedNavPrincipal[item.trade_date].nav + item.nav,
          principal: aggregatedNavPrincipal[item.trade_date].principal + item.principal,
        };
      } else {
        aggregatedNavPrincipal[item.trade_date] = {
          nav: item.nav,
          principal: item.principal,
        };
      }
    }
  });

  const simpleReturns = {"Simple Returns": []};
  const navVsPrincipal = {"NAV": [], "Principal": []};
  Object.entries(aggregatedNavPrincipal).forEach(([date, item]) => {
    simpleReturns["Simple Returns"].push({
      trade_date: date,
      simple_returns: (item.nav - item.principal) / item.principal,
    });

    navVsPrincipal["NAV"].push({
      trade_date: date,
      value: item.nav,
    });
    navVsPrincipal["Principal"].push({
      trade_date: date,
      value: item.principal,
    });
  })

  const latestDate = portfolio.at(-1).trade_date
  const latestSymbols = portfolio.filter(x => x.trade_date === latestDate)
  const total = Object
    .values(latestSymbols)
    .map(x => x.nav)
    .reduce((previous, current) => previous + current)
  const pieChartData = latestSymbols.map((item) => {
    return {
      ...item,
      category: item.symbol,
      percent: (item.nav/total),
    };
  });


  const latestNAV = navVsPrincipal["NAV"]
    .at(-1)
    .value
  const latestPrincipal = navVsPrincipal["Principal"]
    .at(-1)
    .value
  const latestSimpleReturns = (latestNAV - latestPrincipal) / latestPrincipal

  return (
    <EuiPageTemplate
      restrictWidth={false}
      pageSideBar={<SideBar/>}
      pageHeader={{
        iconType: 'logoElastic',
        pageTitle: 'Investments Summary',
      }}
    >
      <EuiFlexGroup>
        <EuiFlexItem grow={6}>
          <EuiBasicTable
            items={pieChartData}
            columns={columns}
          />
        </EuiFlexItem>
        <EuiFlexItem grow={4}>
          <PieChart
            data={pieChartData}
          />
        </EuiFlexItem>
        <EuiFlexItem grow={10}>
          <EuiFlexGroup>
            <EuiFlexItem grow={5}>
              <Stat
                title="Principal"
                value={formatMoney(latestPrincipal)}
                colour="default"
              />
            </EuiFlexItem>
            <EuiFlexItem grow={5}>
              <Stat
                title="NAV"
                value={formatMoney(latestNAV)}
                colour={latestNAV > latestPrincipal ? "success" : "danger"}
              />
            </EuiFlexItem>
          </EuiFlexGroup>
          <EuiSpacer size="m" />
          <EuiFlexGroup>
            <EuiFlexItem grow={5}/>
            <EuiFlexItem grow={5}>
              <Stat
                title="Simple Returns"
                value={formatPercent(latestSimpleReturns)}
                colour={latestSimpleReturns >= 0 ? "success" : "danger"}
              />
            </EuiFlexItem>
          </EuiFlexGroup>
        </EuiFlexItem>
      </EuiFlexGroup>
      <EuiHorizontalRule />
      <EuiSuperSelect
        options={timeOptions}
        valueOfSelected={queryPeriodInMonths}
        onChange={(v) => dispatch(updateQueryPeriodInMonths(v))}
        compressed={true}
        fullWidth={false}
      />
      <EuiSpacer size="m" />
      <EuiFlexGroup>
        <EuiFlexItem>
          <GenericChart
            allSymbols={allSymbols}
            itemKey="simple_returns"
            title="Simple Returns"
            formatCallback={formatPercent}
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <GenericChart
            allSymbols={allSymbols}
            itemKey="nav"
            title="NAV"
            formatCallback={x => formatMoneyGraph(x, 2)}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
      <EuiFlexGroup>
        <EuiFlexItem>
          <GenericChart
            allSymbols={simpleReturns}
            itemKey="simple_returns"
            title="Portfolio Simple Returns"
            formatCallback={formatPercent}
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <GenericChart
            allSymbols={navVsPrincipal}
            itemKey="value"
            title="Portfolio NAV"
            formatCallback={x => formatMoneyGraph(x, 2)}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate>
  );
}

function PieChart({data}) {
  return (
    <Chart size={{height: "250px"}}>
      <Partition
        data={data}
        valueAccessor={d => Number(d.percent)}
        valueFormatter={() => ''} // Hide the slice value if data values are already in percentages
        layers={[
          {
            groupByRollup: d => d.category,
            shape: {
              fillColor: d => EUI_CHARTS_THEME_LIGHT.theme.colors.vizColors[d.sortIndex],
            },
          },
        ]}
        config={{
          ...EUI_CHARTS_THEME_LIGHT.partition,
          clockwiseSectors: false, // For correct slice order
        }}
      />
    </Chart>
  );
}

function GenericChart({allSymbols, itemKey, title, formatCallback}) {
  let minValue = Number.MAX_VALUE;
  let maxValue = Number.NEGATIVE_INFINITY;
  Object.values(allSymbols).forEach((series) => {
    series.forEach((item) => {
      if (item[itemKey] < minValue) {
        minValue = item[itemKey]
      }

      if (item[itemKey] > maxValue) {
        maxValue = item[itemKey]
      }
    });
  });

  const minValueWithBuffer = minValue - (((minValue + maxValue) / 2) * 0.05)

  return (
    <Chart size={{height: "300px"}}>
      <Settings
        theme={EUI_CHARTS_THEME_LIGHT.theme}
        showLegend={true}
        legendPosition="right"
        tooltip={{
        }}
      />
      {
        Object.entries(allSymbols).map(([symbol, series]) => (
          <AreaSeries
            id={`${title}-${symbol}`}
            key={`${title}-${symbol}`}
            name={symbol}
            data={series.map((s => [moment(s.trade_date).valueOf(), s[itemKey]]))}
            xScaleType="time"
            xAccessor={0}
            yAccessors={[1]}
          />
        ))
      }
      <Axis
        id={`${title}-top`}
        title={title}
        position={Position.Top}
        tickFormat={() => ""}
        showGridLines={false}
      />
      <Axis
        id={`${title}-bottom`}
        position={Position.Bottom}
        tickFormat={timeFormatter(niceTimeFormatByDay(31))}
        showGridLines={false}
      />
      <Axis
        id="left-axis"
        position={Position.Left}
        tickFormat={formatCallback}
        showGridLines={false}
        domain={{
          min: minValueWithBuffer,
        }}
      />
    </Chart>
  );
}
