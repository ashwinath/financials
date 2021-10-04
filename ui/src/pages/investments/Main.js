import React from 'react';

import {
  EuiPageTemplate,
} from '@elastic/eui';

import { SideBar } from "../../components";
import { LoadingPage } from "../";
import { useLoginHook } from "../../hooks/login";

export function InvestmentsMainPage() {
  const status = useLoginHook();

  if (status === "loading") {
    return <LoadingPage/>;
  }

  return (
    <EuiPageTemplate
      pageSideBar={<SideBar/>}
      pageHeader={{
        iconType: 'logoElastic',
        pageTitle: 'Investments',
      }}
    >
      <h1>hello world investments</h1>
    </EuiPageTemplate>
  );
}
