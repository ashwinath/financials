import React from 'react';

import { EuiSideNav, EuiIcon } from '@elastic/eui';
import { useLocation, useHistory } from 'react-router-dom';
import { useSelector } from 'react-redux';

const HOME_PAGE = "Home";
const INVESTMENTS_SUMMARY = "Summary"
const INVESTMENTS_TRADES = "Trades"

const PATH_MAPPING = {
  [HOME_PAGE]: "/",
  [INVESTMENTS_SUMMARY]: "/investments",
  [INVESTMENTS_TRADES]: "/investments/trades",
}

export function SideBar() {
  const history = useHistory();
  const { pathname } = useLocation();
  const { isLoggedIn } = useSelector((state) => state.login)

  if (!isLoggedIn) {
    history.push("/login");
  }

  let i = 0;
  const sideNav = [
    {
      name: 'Investments',
      icon: <EuiIcon type="logoElasticsearch" />,
      id: `${i++}`,
      items: [
        {
          name: INVESTMENTS_SUMMARY,
          id: INVESTMENTS_SUMMARY,
          onClick: () => history.push(PATH_MAPPING[INVESTMENTS_SUMMARY]),
          isSelected: pathname === PATH_MAPPING[INVESTMENTS_SUMMARY],
        },
        {
          name: INVESTMENTS_TRADES,
          id: INVESTMENTS_TRADES,
          onClick: () => history.push(PATH_MAPPING[INVESTMENTS_TRADES]),
          isSelected: pathname === PATH_MAPPING[INVESTMENTS_TRADES],
        },
      ],
    },
  ];

  return (
		<EuiSideNav
			mobileTitle="Nav Items"
			items={sideNav}
		/>
  );
}
