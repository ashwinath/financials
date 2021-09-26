import React from 'react';

import {
  useHistory,
} from "react-router-dom";

import {
  EuiPageTemplate,
} from '@elastic/eui';

import { 
  useDispatch,
  useSelector,
} from 'react-redux';

import { querySessionAsync } from "../redux/mainPageSlice";
import Sidebar from "./SideBar";
import LoadingPage from "./LoadingPage";

export default function HomePage() {
  const history = useHistory();
  const dispatch = useDispatch();

  const { isLoggedIn, status, triedLoggingIn } = useSelector((state) => state.mainPage)
  if (!triedLoggingIn && status === "idle") {
    dispatch(querySessionAsync())
  }

  if (!isLoggedIn && triedLoggingIn) {
    history.push("/login");
  }

  if (status === "loading") {
    return <LoadingPage/>;
  }

  return (
    <EuiPageTemplate
      pageSideBar={<Sidebar/>}
      pageHeader={{
        iconType: 'logoElastic',
        pageTitle: 'Financials',
      }}
    >
      <h1>hello world</h1>
    </EuiPageTemplate>
  );
}
