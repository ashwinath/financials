import {
  EuiCollapsibleNavGroup,
} from '@elastic/eui';

function BrandBar() {
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

export default BrandBar;
