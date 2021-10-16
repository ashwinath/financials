import React from 'react';

import {
  EuiPageTemplate,
} from '@elastic/eui';

import { SideBar } from "../components";

export function HomePage() {
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
