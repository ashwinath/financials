import React from 'react';

import {
  EuiPageTemplate,
} from '@elastic/eui';

import { SideBar } from "../../components";
import { LoadingPage } from "../";
import { useLoginHook } from "../../hooks/login";
import { queryTrades } from '../../redux/investmentsSlice';
import { useDispatch, useSelector } from 'react-redux';
import { ErrorBar } from "../../components"

export function InvestmentsTradesPage() {
  const loginStatus = useLoginHook();

  const dispatch = useDispatch();
  const investmentsState = useSelector((state) => state.investments);
  const {
    page,
    pageSize,
    orderBy,
    order,
    status,
    payload,
    errorMessage,
  } = investmentsState;

  if (status === "idle" && payload === null) {
    dispatch(queryTrades({page, pageSize, orderBy, order}));
  }

  if (loginStatus === "loading" || status === "loading") {
    return <LoadingPage/>;
  }

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
          pageTitle: 'Investments Trades',
        }}
      >
        <h1>hello world investments trades page</h1>
      </EuiPageTemplate>
    </>
  );
}
