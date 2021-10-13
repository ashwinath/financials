import React from 'react';

import {
  EuiPageTemplate,
} from '@elastic/eui';

import { SideBar } from "../components";
import { LoadingPage } from ".";
import { useLoginHook } from "../hooks";

export function HomePage() {
  const status = useLoginHook();

  if (status === "loading") {
    return <LoadingPage/>;
  }

  return (
    <EuiPageTemplate
      pageSideBar={<SideBar/>}
      pageHeader={{
        iconType: 'logoElastic',
        pageTitle: 'Financials',
      }}
    >
      <h1>hello world</h1>
    </EuiPageTemplate>
  );
}
