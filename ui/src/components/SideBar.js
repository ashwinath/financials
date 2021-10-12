import React from 'react';

import { htmlIdGenerator, EuiSideNav, EuiIcon } from '@elastic/eui';
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

  const sideNav = [
    {
      name: 'Investments',
      icon: <EuiIcon type="logoElasticsearch" />,
      id: `${htmlIdGenerator()()}`,
      href: PATH_MAPPING[INVESTMENTS_SUMMARY],
      onClick: (e) => {
        e.preventDefault();
        history.push(PATH_MAPPING[INVESTMENTS_SUMMARY])
      },
      isSelected: pathname === PATH_MAPPING[INVESTMENTS_SUMMARY],
      items: [
        {
          name: INVESTMENTS_TRADES,
          id: INVESTMENTS_TRADES,
          href: PATH_MAPPING[INVESTMENTS_TRADES],
          onClick: (e) => {
            e.preventDefault();
            history.push(PATH_MAPPING[INVESTMENTS_TRADES])
          },
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
