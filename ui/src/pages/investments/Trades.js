import React from 'react';

import {
  formatDate,
  EuiPageTemplate,
  EuiBasicTable,
} from '@elastic/eui';

import { SideBar } from "../../components";
import { LoadingPage } from "../";
import { useLoginHook } from "../../hooks/login";
import { queryTrades, updateTableInfo, resetShouldReload, setInitialState } from '../../redux/investmentsSlice';
import { useDispatch, useSelector } from 'react-redux';
import { ErrorBar } from "../../components";
import { capitaliseFirstLetter, capitaliseAll, formatMoney } from "../../utils";
import { useHistory, useLocation } from "react-router-dom";

export function InvestmentsTradesPage() {
  const loginStatus = useLoginHook();
  const history = useHistory();
  const dispatch = useDispatch();
  const location = useLocation();
  const investmentsState = useSelector((state) => state.investments);

  const {
    page,
    pageSize,
    orderBy,
    order,
    status,
    payload,
    errorMessage,
    shouldReload,
    init,
  } = investmentsState;

  if (!init) {
    dispatch(setInitialState(new URLSearchParams(location.search)));
  }

  if (init && status === "idle" && (payload === null || shouldReload)) {
    dispatch(queryTrades({page, pageSize, orderBy, order}));
  }

  if (loginStatus === "loading") {
    return <LoadingPage/>;
  }

  if (shouldReload) {
    dispatch(resetShouldReload());
    history.push({
      pathname: "/investments/trades",
      search: `?page=${page}&page_size=${pageSize}&order=${order}&order_by=${orderBy}`,
    });
  }

  let results = [];
  if (payload && payload.data) {
    results = payload.data.results.map((data) => {
      return {
        ...data,
        total: data.price_each * data.quantity,
      };
    });
  }

  const columns = [
    {
      field: 'date_purchased',
      name: 'Date',
      dataType: 'date',
      render: (date) => formatDate(date, 'dobLong'),
      sortable: true,
    },
    {
      field: 'symbol',
      name: 'Symbol',
      truncateText: true,
      render: (field) => capitaliseAll(field),
    },
    {
      field: 'trade_type',
      name: 'Trade Type',
      truncateText: true,
      render: (field) => capitaliseFirstLetter(field),
    },
    {
      field: 'price_each',
      name: 'Price',
      truncateText: true,
      render: (field) => formatMoney(field),
    },
    {
      field: 'quantity',
      name: 'Quantity',
      truncateText: true,
    },
    {
      field: 'total',
      name: 'Total',
      truncateText: true,
      render: (field) => formatMoney(field),
    }
  ];

  const pagination = {
    pageIndex: payload && payload.data ? payload.data.paging.page - 1 : 0,
    pageSize: pageSize,
    totalItemCount: payload && payload.data ? payload.data.paging.total : 0,
    pageSizeOptions: [10, 20],
    hidePerPageOptions: false,
  };

  const sorting = {
    sort: {
      field: order,
      direction: orderBy,
    },
  };

  return (
    <>
      <ErrorBar 
        title="Sorry, there was an error retrieving your trades."
        errorMessage={errorMessage}
      />

      <EuiPageTemplate
        pageSideBar={<SideBar/>}
        pageHeader={{
          iconType: 'logoElastic',
          pageTitle: 'Trades',
        }}
      >
        <EuiBasicTable
          items={results}
          columns={columns}
          pagination={pagination}
          onChange={(value) => dispatch(updateTableInfo(value))}
          loading={(loginStatus === "loading" || status === "loading")}
          sorting={sorting}
        />
      </EuiPageTemplate>
    </>
  );
}
