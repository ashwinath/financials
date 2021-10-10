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
} from '@elastic/charts';

import { EUI_CHARTS_THEME_LIGHT } from '@elastic/eui/dist/eui_charts_theme';

import { SideBar } from "../../components";
import { LoadingPage } from "../";
import { useLoginHook } from "../../hooks/login";
import { useDispatch, useSelector } from 'react-redux';

import {
  queryPortfolio,
} from '../../redux/investmentsSlice';

import { getDateFromPeriod, formatMoney } from "../../utils";

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

  if (status === "loading" || portfolioLoading) {
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


  return (
    <EuiPageTemplate
      restrictWidth={false}
      pageSideBar={<SideBar/>}
      pageHeader={{
        iconType: 'logoElastic',
        pageTitle: 'Investments',
      }}
    >
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
            formatCallback={formatMoney}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate>
  );
}

function GenericChart({allSymbols, itemKey, title, formatCallback}) {
  return (
    <Chart size={{height: "30vh"}}>
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
        id="bottom-axis"
        position="bottom"
        tickFormat={timeFormatter(niceTimeFormatByDay(31))}
        showGridLines={false}
      />
      <Axis
        id="left-axis"
        position="left"
        tickFormat={formatCallback}
        showGridLines
      />
    </Chart>
  );
}
