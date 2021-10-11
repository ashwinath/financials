import React from 'react';
import {
  EuiPanel,
  EuiStat,
  EuiIcon,
} from '@elastic/eui';

export function Stat({title, value, colour}) {
  return (
    <EuiPanel>
      <EuiStat
        title={value}
        description={title}
        textAlign="right"
        titleColor={colour}
      >
        <EuiIcon type="empty" />
      </EuiStat>
    </EuiPanel>
  );
}
