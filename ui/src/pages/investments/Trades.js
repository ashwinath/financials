import React from 'react';

import {
  formatDate,
  EuiPageTemplate,
  EuiBasicTable,
  EuiTextColor,
  EuiButton,
  EuiSpacer,
  EuiFlexItem,
  EuiFlexGroup,
} from '@elastic/eui';

import { SideBar, AddTradeModal, AddBulkTradeModal } from "../../components";
import {
  queryTrades,
  updateTableInfo,
  resetShouldReload,
  setInitialState,
  toggleIsAddTradeModalOpen,
  toggleIsAddBulkTradeModalOpen,
  setSelectedItems,
  deleteTrades,
} from '../../redux/investmentsSlice';
import { useDispatch, useSelector } from 'react-redux';
import { ErrorBar, SuccessBar } from "../../components";
import { capitaliseFirstLetter, capitaliseAll, formatMoney } from "../../utils";
import { useHistory, useLocation } from "react-router-dom";

export function InvestmentsTradesPage() {
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
    isAddTradeModalOpen,
    isAddBulkTradeModalOpen,
    submitSuccess,
    selectedItems,
    successMessage,
  } = investmentsState;

  if (!init) {
    dispatch(setInitialState(new URLSearchParams(location.search)));
  }

  if (init && status === "idle" && (payload === null || shouldReload)) {
    dispatch(queryTrades({page, pageSize, orderBy, order}));
  }

  if (shouldReload) {
    dispatch(resetShouldReload());
    history.push({
      pathname: "/investments/trades",
      search: `?page=${page}&page_size=${pageSize}&order_by=${orderBy}&order=${order}`,
    });
  }

  let results = [];
  if (payload && payload.results) {
    results = payload.results.map((data) => {
      return {
        ...data,
        total: data.price_each * data.quantity,
      };
    });
  }

  const renderTradeType = (type) => {
    const color = type === "buy" ? 'success' : 'danger';
    return <EuiTextColor color={color}>{capitaliseFirstLetter(type)}</EuiTextColor>;
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
      render: (field) => renderTradeType(field),
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
    pageIndex: payload && payload.paging ? payload.paging.page - 1 : 0,
    pageSize: pageSize,
    totalItemCount: payload && payload.paging ? payload.paging.total : 0,
    pageSizeOptions: [20, 40],
    hidePerPageOptions: false,
  };

  const sorting = {
    sort: {
      field: orderBy,
      direction: order,
    },
  };

  const selection = {
    onSelectionChange: (items) => dispatch(setSelectedItems(items)),
    initialSelected: selectedItems,
  };

  const renderDeleteButton = () => {
    if (selectedItems.length === 0) {
      return null;
    }

    return (
      <EuiButton
        color="danger"
        size="s"
        iconType="trash"
        onClick={() => dispatch(deleteTrades(selectedItems))}
      >
         Delete {selectedItems.length} trades
      </EuiButton>
    );
  }

  return (
    <>
      <ErrorBar 
        title="Sorry, there was an error."
        errorMessage={errorMessage}
      />
      <SuccessBar 
        title={successMessage}
        message={submitSuccess === "success" ? "We did it!" : null}
      />

      <EuiPageTemplate
        pageSideBar={<SideBar/>}
        pageHeader={{
          iconType: 'logoElastic',
          pageTitle: 'Trades',
        }}
      >
        <EuiFlexGroup responsive={false} wrap gutterSize="s" alignItems="center">
          <EuiFlexItem grow={false} alignItems="center">
            <EuiButton
              size="s"
              onClick={() => dispatch(toggleIsAddTradeModalOpen())}
            >
              Add trade
            </EuiButton>
          </EuiFlexItem>
          <EuiFlexItem grow={false}>
            <EuiButton
              size="s"
              onClick={() => dispatch(toggleIsAddBulkTradeModalOpen())}
            >
              Add bulk trade
            </EuiButton>
          </EuiFlexItem>
          <EuiFlexItem/>
          {renderDeleteButton()}
        </EuiFlexGroup>
        <EuiSpacer size="s" />
        <EuiBasicTable
          items={results}
          itemId="id"
          isSelectable={true}
          selection={selection}
          columns={columns}
          pagination={pagination}
          onChange={(value) => dispatch(updateTableInfo(value))}
          loading={status === "loading"}
          sorting={sorting}
        />
      </EuiPageTemplate>
      {isAddTradeModalOpen ? <AddTradeModal/> : null}
      {isAddBulkTradeModalOpen ? <AddBulkTradeModal/> : null}
    </>
  );
}
