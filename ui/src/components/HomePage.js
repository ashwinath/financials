import React from 'react';

import {
  EuiPageTemplate,
} from '@elastic/eui';

import Sidebar from "./SideBar";
import LoadingPage from "./LoadingPage";
import { useLoginHook } from "../hooks/login";

export default function HomePage() {
  const status = useLoginHook();

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
