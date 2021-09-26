import React from 'react';

import {
  EuiPageTemplate,
  EuiEmptyPrompt,
  EuiLoadingChart,
} from '@elastic/eui';

export default function Loading() {
  return (
    <EuiPageTemplate
      template="centeredBody"
      pageContentProps={{ paddingSize: 'l' }}
      minHeight="80vh"
    >
      <EuiEmptyPrompt
        title={<span>Personalising your experience</span>}
        body={<EuiLoadingChart size="xl"/>}
        titleSize="m"
      />
    </EuiPageTemplate>
  );
}
