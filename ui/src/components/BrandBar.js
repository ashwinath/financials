import { EuiCollapsibleNavGroup } from '@elastic/eui';

export function BrandBar() {
  return (
    <EuiCollapsibleNavGroup
      title="Financials"
      iconType="logoMetrics"
      iconSize="l"
      titleSize="s"
      isCollapsible={false}
      initialIsOpen={false}
      background="dark"
    />
  );
}
