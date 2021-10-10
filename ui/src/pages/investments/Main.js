import React from 'react';
import moment from 'moment';

import {
  EuiPageTemplate,
  EuiFlexGroup,
  EuiFlexItem,
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

import { SideBar } from "../../components";
import { LoadingPage } from "../";
import { useLoginHook } from "../../hooks/login";
import { useDispatch, useSelector } from 'react-redux';

import {
  queryPortfolio,
} from '../../redux/investmentsSlice';

import { getDateFromPeriod, formatMoneyGraph } from "../../utils";

export function InvestmentsMainPage() {
  const status = useLoginHook();
  const dispatch = useDispatch();
  const {
    queryPeriodInMonths,
    portfolio,
    portfolioLoaded,
    portfolioLoading,
  } = useSelector((state) => state.investments);

  const startDate = getDateFromPeriod(queryPeriodInMonths)
  if (!portfolioLoading && !portfolioLoaded) {
    dispatch(queryPortfolio(startDate))
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
      if (aggregatedNavPrincipal[item.trade_date] in aggregatedNavPrincipal) {
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
  const total = Object.values(latestSymbols).map(x => x.nav).reduce((previous, current) => previous + current)
  const pieChartData = latestSymbols.map((item) => {
    return {
      category: item.symbol,
      percent: (item.nav/total),
    };
  });

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
        <EuiFlexItem>
          <PieChart
            data={pieChartData}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
      <EuiFlexGroup>
        <EuiFlexItem>
          <GenericChart
            allSymbols={allSymbols}
            itemKey="simple_returns"
            title="Simple Returns"
            formatCallback={d => `${Number(d * 100).toFixed(2)}%`}
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
            title="Simple Returns"
            formatCallback={d => `${Number(d * 100).toFixed(2)}%`}
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <GenericChart
            allSymbols={navVsPrincipal}
            itemKey="value"
            title="Simple Returns"
            formatCallback={x => formatMoneyGraph(x, 2)}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate>
  );
}

function PieChart({data}) {
  return (
    <Chart size={{height: "25vh"}}>
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
  return (
    <Chart size={{height: "25vh"}}>
      <Settings
        theme={EUI_CHARTS_THEME_LIGHT.theme}
        showLegend={true}
        legendPosition="right"
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
        id={title}
        title={title}
        position={Position.Bottom}
        tickFormat={timeFormatter(niceTimeFormatByDay(31))}
        showGridLines={false}
      />
      <Axis
        id="left-axis"
        position={Position.Right}
        tickFormat={formatCallback}
        showGridLines
      />
    </Chart>
  );
}
